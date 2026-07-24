-- name: GetDiaries :many
SELECT diaries.*, alters.avatar FROM diaries join alters on diaries.alter_id = alters.id 
WHERE diaries.user_id = $1;