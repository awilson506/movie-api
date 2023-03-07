package queries

// GetProductionCompanyDetailsSQL - query to get company details by year and company
const GetProductionCompanyDetailsSQL = `
SELECT mpc.id AS id, mpc.name, SUM(budget) AS budget, release_year AS year, SUM(revenue) AS revenue
FROM movies m
JOIN movies_to_production_companies mtpc
	ON m.id = mtpc.movie_id
JOIN production_companies mpc
	ON mpc.id = mtpc.production_company_id
	$CONDITIONALS
GROUP BY mpc.id, m.release_year
ORDER BY mpc.id
LIMIT $PAGE
`

// GetGenreDetailsSQL - query to get genre details by year
const GetGenreDetailsSQL = `
SELECT name, release_year AS year, genre_id, MAX(count) AS count 
FROM (
SELECT g.id AS genre_id, g.name, COUNT(g.id) AS count, SUM(budget) AS budget, release_year
FROM movies m
JOIN movies_to_genres mtg
	ON m.id = mtg.movie_id
JOIN genres g
	ON g.id = mtg.genre_id
GROUP BY g.id, m.release_year
ORDER BY count DESC
)
$CONDITIONALS
GROUP BY release_year
LIMIT $PAGE
`
