-- name: CreateAlter :one
INSERT INTO alters (id, created_at, updated_at, name, pronouns, age, alter_role, description, colour, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING  *;
