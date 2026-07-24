-- name: GetDiaryEntry :one
SELECT * FROM diary_entries WHERE id = $1 and user_id = $2;
