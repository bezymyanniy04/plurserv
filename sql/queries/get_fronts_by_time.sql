-- name: GetFrontsByTime :many
SELECT *
FROM fronts JOIN alters ON fronts.alter_id = alters.id
WHERE alters.user_id = $3 and ((fronts.started_at between $1 and $2) or (fronts.ended_at between $1 and $2)
or (fronts.started_at < $1 and(fronts.ended_at IS NULL or fronts.ended_at > $2))) ;
