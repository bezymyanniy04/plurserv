-- name: NewDiaryEntry :one
INSERT INTO diary_entries (id, diary_id, name, date, text, user_id)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    (NOW() at time zone 'utc'),
    $3,
    $4
)
RETURNING *;