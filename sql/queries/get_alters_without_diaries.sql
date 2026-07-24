-- name: GetAltersWithoutDiaries :many
SELECT alters.* FROM alters
WHERE alters.user_id = $1 and alters.id not in(SELECT alter_id FROM diaries)
ORDER BY alters.created_at ASC
;
