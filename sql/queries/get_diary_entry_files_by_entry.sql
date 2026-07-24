-- name: GetDiaryEntryFilesByEntry :many
SELECT * FROM diary_files WHERE entry_id = $1 and user_id = $2;
