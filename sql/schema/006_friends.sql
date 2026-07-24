-- +goose Up
CREATE TABLE friends(
id UUID PRIMARY KEY,
request_id UUID NOT NULL,
user_id  UUID NOT NULL,
friend_id  UUID NOT NULL,
CONSTRAINT fk_request_id FOREIGN KEY (request_id) REFERENCES requests(id) ON DELETE CASCADE,
CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
CONSTRAINT fk_friend_id FOREIGN KEY (friend_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE friends;