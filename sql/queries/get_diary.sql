-- name: GetDiary :one
SELECT * FROM diaries 
WHERE id = $1 and user_id = $2;