-- name: NewFriends :exec
INSERT INTO friends (id, request_id, user_id, friend_id)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3
);
