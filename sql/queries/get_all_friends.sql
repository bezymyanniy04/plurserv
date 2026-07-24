-- name: GetAllFriends :many
SELECT * FROM users JOIN friends
ON users.id = friends.friend_id
WHERE friends.user_id = $1;