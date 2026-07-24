-- name: DeleteFriendRequest :exec
DELETE FROM requests WHERE id = $1;