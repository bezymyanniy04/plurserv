-- name: NewDiary :one
INSERT INTO diaries (id, alter_id, bg_colour, bg_colour2, block_colour, text_colour, user_id, name)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;
