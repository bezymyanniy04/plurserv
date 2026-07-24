-- name: EditForNewbies :one
UPDATE for_newbies
SET text = $1
WHERE user_id = $2
RETURNING *;