
-- name: GetRouteByPath :one
SELECT * FROM routes
WHERE path = $1 LIMIT 1;

