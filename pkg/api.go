package api

import (
	"github.com/awilson506/movie-api/database"
)

const PRODUCTION_COMPANY_ID = "production_company_id"
const YEAR = "release_year"

// Client - hold the db connector or cache client if we had one
type Client struct {
	db database.DbConfig
}

// New - get a new instance of Client
func New() *Client {
	return &Client{
		db: *database.NewDbConnection(),
	}
}

// GetProductionCompanyDetails - map out parameters and get the company details
func (c *Client) GetProductionCompanyDetails(productionCompanyId *string, year *string, pageId *int) []database.ProductionCompanyDetails {
	conditionals := map[string]*string{
		PRODUCTION_COMPANY_ID: productionCompanyId,
		YEAR:                  year,
	}
	productionCompanyDetails := c.db.GetProductionCompanyDetails(conditionals, pageId)
	return productionCompanyDetails
}

// GetGenreDetails - map out parameters and get the genre details
func (c *Client) GetGenreDetails(year *string, pageId *int) []database.GenreDetails {
	conditionals := map[string]*string{
		YEAR: year,
	}
	genreDetails := c.db.GetGenreDetails(conditionals, pageId)
	return genreDetails
}
