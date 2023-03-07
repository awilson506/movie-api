package database

import (
	"fmt"
	"log"
	"strings"

	"github.com/awilson506/movie-api/database/queries"
)

type ProductionCompanyDetails struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Budget  int    `json:"budget"`
	Year    int    `json:"year"`
	Revenue int    `json:"revenue"`
}

type GenreDetails struct {
	Id    int    `json:"id"`
	Year  int    `json:"year"`
	Name  string `json:"name"`
	Count int    `json:"-"`
}

// GetProductionCompanyDetails
func (db *DbConfig) GetProductionCompanyDetails(conditionals map[string]*string, pageId *int) []ProductionCompanyDetails {
	productionCompanyDetails := []ProductionCompanyDetails{}
	parsedQuery := appendConditionals(queries.GetProductionCompanyDetailsSQL, conditionals)
	parsedQuery = setPagination(parsedQuery, pageId)

	rows, err := db.Connection.Query(parsedQuery)

	if err != nil {
		log.Fatalf("read failed: %s", err)
	}

	for rows.Next() {
		var pcd ProductionCompanyDetails
		err = rows.Scan(&pcd.Id, &pcd.Name, &pcd.Budget, &pcd.Year, &pcd.Revenue)
		if err != nil {
			log.Fatalf("Scan: %v", err)
		}
		productionCompanyDetails = append(productionCompanyDetails, pcd)
	}
	return productionCompanyDetails
}

// GetGenreDetails
func (db *DbConfig) GetGenreDetails(conditionals map[string]*string, pageId *int) []GenreDetails {
	genreDetails := []GenreDetails{}
	parsedQuery := appendConditionals(queries.GetGenreDetailsSQL, conditionals)
	parsedQuery = setPagination(parsedQuery, pageId)
	rows, err := db.Connection.Query(parsedQuery)

	if err != nil {
		log.Fatalf("read failed: %s", err)
	}

	for rows.Next() {
		var gd GenreDetails
		err = rows.Scan(&gd.Name, &gd.Year, &gd.Id, &gd.Count)
		if err != nil {
			log.Fatalf("Scan: %v", err)
		}
		genreDetails = append(genreDetails, gd)
	}
	return genreDetails
}

func appendConditionals(query string, conditionals map[string]*string) string {
	var conditionalSubQuery string
	cnt := 0

	if len(conditionals) > 0 {
		for i, v := range conditionals {
			if cnt == 0 && *v != "" {
				conditionalSubQuery = "WHERE " + i + " = " + *v
				cnt++
			} else if *v != "" {
				conditionalSubQuery = conditionalSubQuery + " AND " + i + " = " + *v
				cnt++
			}
		}
	}
	return strings.Replace(query, "$CONDITIONALS", conditionalSubQuery, -1)
}

func setPagination(query string, pageId *int) string {
	pageQuery := "0, 100"
	if *pageId > 1 {
		var pageLimit int = 100 * *pageId
		pageQuery = fmt.Sprintf("%d,%d", pageLimit-99, 100)
	}

	return strings.Replace(query, "$PAGE", pageQuery, -1)
}
