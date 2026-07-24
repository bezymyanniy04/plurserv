-- +goose Up
CREATE TABLE diary_files(
id UUID PRIMARY KEY,
created_at TIMESTAMP NOT NULL,
entry_id  UUID NOT NULL,
CONSTRAINT fk_entry_id
FOREIGN KEY (entry_id)
REFERENCES diary_entries(id) ON DELETE CASCADE,
file TEXT NOT NULL,
user_id UUID NOT NULL,
CONSTRAINT fk_user_id
FOREIGN KEY (user_id)
REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE diary_files;

