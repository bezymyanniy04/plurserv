-- name: GetNowFrontByAlter :one
SELECT * FROM fronts WHERE alter_id = $1 and ended_at IS NULL;
