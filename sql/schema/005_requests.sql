-- +goose Up
CREATE TABLE requests(
id UUID PRIMARY KEY,
created_at TIMESTAMP NOT NULL,
expires_at TIMESTAMP default NULL,
sender_id  UUID NOT NULL,
receiver_id  UUID NOT NULL,
CONSTRAINT fk_sender_id FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
CONSTRAINT fk_receiver_id FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE,
answer INTEGER NOT NULL default 0
);

-- +goose Down
DROP TABLE requests;
