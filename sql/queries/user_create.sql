-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, avatar, system_name, theme, font)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING  *;
