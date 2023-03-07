#!/bin/bash

echo -n "Production Company Details get call: "
COMPANY_DETAILS=$(curl -s http://localhost:8080/production-company-details?year=1900)
echo $COMPANY_DETAILS | python -m json.tool


echo "Genre Details get call: "
GENRE_DETAILS=$(curl -s http://localhost:8080/genre-details?year=1928)


echo $GENRE_DETAILS | python -m json.tool
