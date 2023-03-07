package database

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func (db *DbConfig) PreProcessCsvData() {
	dataFilepath, _ := filepath.Abs(db.DbPath + "/../dataset/movies_metadata.csv")

	rawInput, err := ioutil.ReadFile(dataFilepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// remove some garbage paragraph spaces that cause a new line on the movie description field
	output := strings.Replace(string(rawInput), "\n ", "  ", -1)
	newOutput, err := os.Create(db.DbPath + "/../dataset/movies_processed.csv")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, err := newOutput.WriteString(output); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
