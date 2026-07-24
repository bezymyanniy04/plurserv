-- name: NewForNewbies :one
INSERT INTO for_newbies (id, user_id, text)
VALUES (
    gen_random_uuid(),
    $1,
    $2
)
RETURNING *;
