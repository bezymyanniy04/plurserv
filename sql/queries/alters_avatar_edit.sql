-- name: EditAlterAvatar :one
UPDATE alters
SET avatar = $3
WHERE id = $1 and user_id = $2
RETURNING *;