-- name: GetAlter :one
SELECT * FROM alters WHERE id = $1;
