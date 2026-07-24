-- name: EditUserSettings :one
UPDATE users
SET
system_name = $1,
theme = $2,
font = $3  
WHERE id = $4
RETURNING *;