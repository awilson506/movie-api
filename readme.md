# Golang Movie Details API
A small application that can take requests for details about movie productions companies and their associated genres.  This application
also has a backend import service to allow movie data to be imported from an s3 bucket and structured in a relational database.  

## Getting the service setup & running
If you don't already have go installed you can download it and install it from here [Download Golang](https://go.dev/doc/install)

Note this application requires the latest go version `1.19` but supports back to `1.17` if you already have a previous version
installed.  You will have to update the go module to require a later version if you choose not to update.  This can be done by 
running: 
```sh
go mod edit -go=1.MY_VERSION
```

## Download and unpack the movie data
In the root directory run:
```sh
go cmd/download/download_data.go
```

## Create the database and import the csv data
In the root directory run:
```sh
go cmd/import/import_data.go
```

## Start the server
In the root directory run:
```sh
go run cmd/server/main.go
```

## Using the api
This application offers 2 endpoints:

### Get Company Details
Takes in three optional parameters and returns the total budget and revenue per year per production company:
```
page=1
year=1990
production_id=2
```
```
curl -s http://localhost:8080/production-company-details
```
Example Response:
```json
[
    {
        "id": 2,
        "name": "Walt Disney Pictures",
        "budget": 37931000,
        "year": 1990,
        "revenue": 47431461
    }
]
```

### Get Genre Details
Takes in two optional parameters and returns the most popular Genre by year:
```
page=1
year=1990
```
```
curl -s http://localhost:8080/genre-details
```
Example Response:
```json
[
    {
        "id": 18,
        "year": 1990,
        "name": "Drama"
    }
]
```

## Testing
Included in the `tests/` directory is a `./load_test.sh` script for running many get requests 
against the `Production Company Details` endpoint. 

Along with the load test, a postman Library is also available in the [tests/postman](/tests/postman/) directory.

## Database
This application uses a small sqLite database to store the movie data.  The ERD diagram can be found [here](/diagrams/movies_db_erg.png).


## Project Notes

### Database
Ideally this application would use a more powerful database like Postgres with real data types and functionality.  But for the purpose of keeping this simple, I used sqLite for local development and testing.

Initially I did entertain the idea of using a noSql database like MongoDB, but after reviewing the data in the CSV and creating use cases for the application the data proved to be inherently relational in nature. 

### Data integrity and import
Importing the data from the CSV and trying to make assumptions about some fields being "json" and others being strings and ints proved to be tough.  I took the approach of assuming most items were string and text and tried to catch and correct as many json syntax issues on some fields during import by looking at the byte codes and trying to fix patterns as they showed up. 

A note about getting this going in production, hopefully we would be able to control the source of the data better and correct issues before we get to the import process. 

### Application env
In the interest of time and keeping the application setup simple, I opted to not leverage env variables for things like db location, s3 paths and csv file locations.  In production, this application would absolutely leverage environment variables for any dynamic non application specific fields. 

### Caching
In production I would have some sort of write through caching layer in-between the database and api service.  Since the application isn't really write heavy it would be easy to maintain a distributed write through cache layer(redis) and only update it when records get updated or created by the import service.  A diagram of the applications can be found [here](/diagrams/movie_api_system_design.png), I included a caching service in the diagram to show how it could work. 