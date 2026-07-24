-- name: GetIfFriends :one
SELECT * FROM friends WHERE user_id = $1 and friend_id = $2;