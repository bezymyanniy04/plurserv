-- name: DeleteDiaryEntry :exec
DELETE FROM diary_entries WHERE id = $1 and user_id = $2;