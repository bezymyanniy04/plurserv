-- name: GetDiaryEntryFile :one
SELECT * FROM diary_files WHERE id = $1 and user_id = $2;
