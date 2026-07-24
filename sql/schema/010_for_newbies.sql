-- +goose Up
CREATE TABLE for_newbies(
id UUID PRIMARY KEY,
user_id  UUID NOT NULL,
CONSTRAINT fk_user_id
FOREIGN KEY (user_id)
REFERENCES users(id) ON DELETE CASCADE,
text TEXT NOT NULL
);

-- +goose Down
DROP TABLE for_newbies;

