-- name: GetDiaryEntriesByDiary :many
SELECT * FROM diary_entries WHERE diary_id = $1 and user_id = $2;
