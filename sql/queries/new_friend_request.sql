-- name: NewFriendRequest :one
INSERT INTO requests (id, created_at, sender_id, receiver_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    $1,
    $2
)
RETURNING *;