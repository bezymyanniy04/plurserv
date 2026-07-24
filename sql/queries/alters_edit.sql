-- name: EditAlter :one
UPDATE alters
SET name = $3,
pronouns = $4,
age = $5,
alter_role = $6,
description = $7, 
colour = $8
WHERE id = $1 and user_id = $2
RETURNING *;