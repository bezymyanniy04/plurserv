-- name: GetAllNowFronts :many
SELECT * FROM alters
WHERE user_id = $1 and fronting = true and LOWER(name) LIKE $2
ORDER BY created_at ASC;
