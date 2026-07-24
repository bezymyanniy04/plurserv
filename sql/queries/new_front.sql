-- name: NewFront :exec
INSERT INTO fronts (id, started_at, alter_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    $1
);

