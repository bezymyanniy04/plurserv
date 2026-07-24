-- name: GetFriendRequest :one
SELECT * FROM requests WHERE (receiver_id = $1 and sender_id = $2) or (receiver_id = $2 and sender_id = $1);