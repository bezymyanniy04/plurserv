-- name: GetNewFriendRequestsReciever :many

SELECT users.*, requests.id, requests.answer FROM requests
JOIN users ON requests.sender_id=users.id
WHERE requests.receiver_id = $1 and requests.answer = 0;