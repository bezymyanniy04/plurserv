-- name: DeleteDiaryEntryFile :exec
DELETE FROM diary_files WHERE id = $1 and user_id = $2;