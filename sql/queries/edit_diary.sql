-- name: EditDiary :one
UPDATE diaries
SET bg_colour = $1,
bg_colour2 = $2, 
block_colour = $3,
text_colour = $4,
font = $5
WHERE id = $6 and user_id = $7
RETURNING *;