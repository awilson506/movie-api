package database

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
	_ "github.com/mattn/go-sqlite3"
)

// MovieMeta - holds the moview meta from the csv import
type MovieMeta struct {
	Budget              int    `csv:"budget"`
	Genres              string `csv:"genres"`
	Id                  int    `csv:"id"`
	ProductionCompanies string `csv:"production_companies"`
	ReleaseDate         string `csv:"release_date"`
	Revenue             int    `csv:"revenue"`
	Runtime             string `csv:"revenue"`
	Name                string `csv:"title"`
}

// ProductionCompany - holds the production company from the csv import
type ProductionCompany struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// Genres - hold the genre object imported from the csv
type Genres struct {
	Id   int
	Name string
}

// LoadDatabase - load the csv data into the database
func (db *DbConfig) LoadDatabase() {
	// this CSV is really dirty, we need to fix some of the PS breaks in the movie descriptions
	// and later sanitize the "json" that is in some of the fields
	db.PreProcessCsvData()

	f, err := os.Open(db.DbPath + "/../dataset/movies_processed.csv")
	if err != nil {
		log.Fatalf("open failed: %s", err)
	}

	defer f.Close()
	movies := []*MovieMeta{}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.FieldsPerRecord = -1
		return r
	})

	if err := gocsv.UnmarshalFile(f, &movies); err != nil {
		panic(err)
	}

	log.Println("loading tables")
	for _, movieMeta := range movies {
		var year int
		productionCompanyIds := db.createProductionCompanies(movieMeta.ProductionCompanies)
		genreIds := db.createGenres(movieMeta.Genres)

		//TODO: make this process asynchronous with a go routine
		db.createGenreRelationships(movieMeta.Id, genreIds)
		db.createProductionCompanyRelationships(movieMeta.Id, productionCompanyIds)

		statement, err := db.Connection.Prepare("INSERT OR IGNORE INTO movies(id, name, budget, revenue, release_date, release_year) values(?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatalf("insert prepare failed: %s", err)
		}
		parts := strings.Split(movieMeta.ReleaseDate, "-")

		if len(parts[0]) == 4 {
			year, err = strconv.Atoi(parts[0])
			if err != nil {
				log.Fatalf("insert prepare failed: %s", err)
			}
		}
		_, err = statement.Exec(movieMeta.Id, movieMeta.Name, movieMeta.Budget, movieMeta.Revenue, movieMeta.ReleaseDate, year)
		if err != nil {
			log.Fatalf("insert failed(%d): %s", movieMeta.Id, err)
		}
	}
	log.Println("tables loaded")
}

// createProductionCompanies - take in json and insert a company for each object
func (db *DbConfig) createProductionCompanies(companies string) []int {
	var productionCompanies []ProductionCompany
	var productionCompanyIds []int

	json.Unmarshal(convertToJsonString(companies), &productionCompanies)

	statement, err := db.Connection.Prepare("INSERT OR IGNORE INTO production_companies(id, name) values(?, ?)")
	if err != nil {
		log.Fatalf("insert prepare failed: %s", err)
	}

	for _, productionCompany := range productionCompanies {
		_, err = statement.Exec(productionCompany.Id, productionCompany.Name)
		if err != nil {
			log.Fatalf("insert failed(%v): %s", productionCompany.Id, err)
		}
		productionCompanyIds = append(productionCompanyIds, productionCompany.Id)
	}
	return productionCompanyIds
}

// createGenres - take in json and insert a genre for each object
func (db *DbConfig) createGenres(genresString string) []int {
	var genres []Genres
	var genreIds []int

	err := json.Unmarshal(convertToJsonString(genresString), &genres)

	if err != nil {
		fmt.Println("unable to import genre: ", err)
	}

	statement, err := db.Connection.Prepare("INSERT OR IGNORE INTO genres(id, name) values(?, ?)")
	if err != nil {
		log.Fatalf("insert prepare failed: %s", err)
	}

	for _, genre := range genres {
		_, err = statement.Exec(genre.Id, genre.Name)
		if err != nil {
			log.Fatalf("insert failed(%v): %s", genre.Id, err)
		}
		genreIds = append(genreIds, genre.Id)
	}
	return genreIds
}

// createGenreRelationships - update the pivot table to reflect the movie to genre relationships
func (db *DbConfig) createGenreRelationships(movieId int, genres []int) {
	statement, err := db.Connection.Prepare("INSERT OR IGNORE INTO movies_to_genres(movie_id, genre_id) values(?, ?)")
	if err != nil {
		log.Fatalf("insert prepare failed: %s", err)
	}

	for _, genre := range genres {
		_, err = statement.Exec(movieId, genre)
		if err != nil {
			log.Fatalf("insert failed(%v): %s", genre, err)
		}
	}
}

// createProductionCompanyRelationships - update the pivot table to reflect the movie to production company relationships
func (db *DbConfig) createProductionCompanyRelationships(movieId int, productionCompanyIds []int) {
	statement, err := db.Connection.Prepare("INSERT OR IGNORE INTO movies_to_production_companies(movie_id, production_company_id) values(?, ?)")
	if err != nil {
		log.Fatalf("insert prepare failed: %s", err)
	}

	for _, productionCompanyId := range productionCompanyIds {
		_, err = statement.Exec(movieId, productionCompanyId)
		if err != nil {
			log.Fatalf("insert failed(%v): %s", productionCompanyId, err)
		}
	}
}

// convertToJsonString clean up dirty "json" strings and attempt to fix missing and single quotes
func convertToJsonString(dirtyString string) []byte {
	data := []byte(dirtyString)
	cleanString := []byte{}
	var shouldReplace = true
	var foundMiddleQuotes = false

	// this is fairly primitive and more cases likely need to be added to it to catch all the posibilities that can occur
	for character := range data {
		if data[character] == byte(92) && data[character+1] == byte(120) {
			cleanString = append(cleanString, byte(92))
		}
		//escape double quotes within quotes and some backslashes
		if data[character] == byte(34) {
			isAlphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(string(data[character+1]))
			if data[character+1] == byte(39) || data[character+1] == byte(32) || (isAlphanumeric && data[character-2] != byte(58)) {
				cleanString = append(cleanString, byte(92))
				foundMiddleQuotes = true
			}
			//replace single quotes
			cleanString = append(cleanString, data[character])
			if !foundMiddleQuotes {
				shouldReplace = !shouldReplace
			} else {
				foundMiddleQuotes = !foundMiddleQuotes
			}
			continue
		}
		// don't replace single quotes while in the middle of double quotes
		if data[character] == 39 && shouldReplace {
			cleanString = append(cleanString, byte(34))
			continue
		}
		if data[character] == 125 && (data[character+1] != 44) {
			// because gatorade can show up in your broken json, try to recover if we don't see the correct pattern
			cleanString = append(cleanString, byte(125))
			cleanString = append(cleanString, byte(93))
			break
		}
		cleanString = append(cleanString, data[character])
	}
	return cleanString
}
