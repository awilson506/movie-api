package main

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/awilson506/movie-api/pkg/s3"
)

func main() {
	var _, b, _, _ = runtime.Caller(0)
	var basepath = filepath.Dir(b)

	dbFilepath, _ := filepath.Abs(basepath + "/../../dataset/the-movies-dataset.zip")
	// download csv data from s3 bucket
	s3.DownloadFromS3Bucket("com.guild.us-west-2.public-data/project-data", "the-movies-dataset.zip", basepath+"/../../dataset/")

	log.Println("unpacking zip to: dataset/")
	// unzip the data into the dataset directory
	s3.UnzipSource(dbFilepath, basepath+"/../../dataset/")
}
