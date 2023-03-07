package main

import (
	"github.com/awilson506/movie-api/database"
)

func main() {

	database.CreateDatabase()
	db := database.NewDbConnection()

	defer db.Connection.Close()

	// create db tables
	db.CreateTables()
	// load data into the db from the csv
	db.LoadDatabase()
}
