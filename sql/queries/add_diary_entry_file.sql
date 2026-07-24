-- name: NewDiaryEntryFile :one
INSERT INTO diary_files (id, entry_id, file, user_id, created_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    NOW()
)
RETURNING *;
