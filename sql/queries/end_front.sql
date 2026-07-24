-- name: EndFrontFronts :exec
UPDATE fronts
SET ended_at = NOW()
WHERE Id = $1 and ended_at IS NULL;
