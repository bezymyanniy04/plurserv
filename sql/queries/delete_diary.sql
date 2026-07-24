-- name: DeleteDiary :exec
DELETE FROM diaries WHERE id = $1 and user_id = $2;