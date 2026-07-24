-- name: EditDiaryEntry :one
UPDATE diary_entries
SET name = $1, 
text = $2
WHERE id = $3 and user_id = $4
RETURNING *;