-- SQLite
SELECT mpc.id AS id, SUM(budget) AS budget, release_year AS year, SUM(revenue)
FROM movies m
JOIN movies_to_production_companies mtpc
ON m.id = mtpc.movie_id
JOIN production_companies mpc
ON mpc.id = mtpc.production_company_id
--WHERE production_company_id = 2
--AND release_year = 1928
GROUP BY mpc.id, m.release_year
ORDER BY mpc.id
LIMIT 1, 100

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
WHERE release_year = 1970
GROUP BY release_year
LIMIT 0,1000

SELECT mpc.id AS id, SUM(budget) AS budget, release_year AS year, SUM(revenue) AS revenue\nFROM movies m\nJOIN movies_to_production_companies mtpc\n\tON m.id = mtpc.movie_id\nJOIN production_companies mpc\n\tON mpc.id = mtpc.production_company_id\n\tAND year = 1928\nGROUP BY mpc.id, m.release_year\nORDER BY mpc.id\nLIMIT 100\n