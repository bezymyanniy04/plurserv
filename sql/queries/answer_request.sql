-- name: AnswerFriendRequest :one
UPDATE requests
SET answer = $1
WHERE receiver_id = $2 and id = $3
RETURNING *;