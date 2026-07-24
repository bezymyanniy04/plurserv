-- name: EditUserInfo :one
UPDATE users
SET email = $1, 
hashed_password = $2,

system_name = $3,
theme = $4,
font = $5  
WHERE id = $6
RETURNING *;