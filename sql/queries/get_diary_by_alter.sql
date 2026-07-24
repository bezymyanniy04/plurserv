-- name: GetDiaryByALter :one
SELECT * FROM diaries 
WHERE alter_id = $1 and user_id = $2;