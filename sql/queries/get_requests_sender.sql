-- name: GetNewFriendRequestsSender :many
SELECT users.*, requests.id, requests.answer FROM requests
JOIN users ON requests.receiver_id=users.id
WHERE requests.sender_id = $1 and requests.answer = 0;