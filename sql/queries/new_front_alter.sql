-- name: NewFrontAlter :one
UPDATE alters
SET fronting = true, fronting_since = (NOW() at time zone 'utc')
WHERE id = $1
RETURNING *;