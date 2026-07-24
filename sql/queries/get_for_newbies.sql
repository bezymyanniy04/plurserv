-- name: GetForNewbies :one
SELECT * FROM for_newbies WHERE user_id = $1;
