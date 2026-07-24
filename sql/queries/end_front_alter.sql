-- name: EndFrontAlter :one
UPDATE alters
SET fronting = false
WHERE id = $1
RETURNING *;