-- name: DeleteAlter :exec
DELETE FROM alters
WHERE id = $1;
