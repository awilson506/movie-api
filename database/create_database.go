package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
)

type DbConfig struct {
	Connection *sql.DB
	DbPath     string
}

func NewDbConnection() *DbConfig {
	// not the best way to get the cd, but this allows the command to be run from any directory(debugger or command line..)
	var _, b, _, _ = runtime.Caller(0)
	var basepath = filepath.Dir(b)

	// file/db paths hsould be in the env but for the sake of cross os usability we will hard code them
	dbFilepath, _ := filepath.Abs(basepath + "/movies.db")

	sqliteDatabase, _ := sql.Open("sqlite3", dbFilepath)

	return &DbConfig{
		Connection: sqliteDatabase,
		DbPath:     basepath,
	}
}

func CreateDatabase() {
	var _, b, _, _ = runtime.Caller(0)
	var basepath = filepath.Dir(b)

	dbFilepath, _ := filepath.Abs(basepath + "/movies.db")
	os.Remove(dbFilepath)

	log.Println("creating movies.db...")
	file, err := os.Create(dbFilepath)
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("movies.db created")
}

func (db *DbConfig) CreateTables() {
	createProductionCompaniesTableSQL := `CREATE TABLE production_companies (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"name" TEXT
	  );`

	createGenresTableSQL := `CREATE TABLE genres (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT
		);`

	createMoviesTableSQL := `CREATE TABLE movies (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT,
		"budget" INTEGER,
		"revenue" INTEGER,
		"release_date" TEXT,
		"release_year" INTEGER
	);`

	createMoviesToProductionCompaniesTableSQL := `CREATE TABLE movies_to_production_companies (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"movie_id" INTEGER,
		"production_company_id" INTEGER,
		FOREIGN KEY("movie_id") REFERENCES movies("id"),
		FOREIGN KEY("production_company_id") REFERENCES production_companies("id")
	);`

	createMoviesTOGenresTableSQL := `CREATE TABLE movies_to_genres (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"movie_id" INTEGER,
		"genre_id" INTEGER,
		FOREIGN KEY("movie_id") REFERENCES movies("id"),
		FOREIGN KEY("genre_id") REFERENCES genres("id")
	);`

	tables := map[int]string{
		0: createProductionCompaniesTableSQL,
		1: createGenresTableSQL,
		2: createMoviesTableSQL,
		3: createMoviesToProductionCompaniesTableSQL,
		4: createMoviesTOGenresTableSQL,
	}

	log.Println("creating tables...")

	for _, createStatementSQL := range tables {
		statement, err := db.Connection.Prepare(createStatementSQL)
		if err != nil {
			log.Fatal(err.Error())
		}
		statement.Exec()
	}
	log.Println("tables created")
}
