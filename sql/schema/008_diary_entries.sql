-- +goose Up
CREATE TABLE diary_entries(
id UUID PRIMARY KEY,
diary_id  UUID NOT NULL,
CONSTRAINT fk_diary_id
FOREIGN KEY (diary_id)
REFERENCES diaries(id) ON DELETE CASCADE,
name TEXT NOT NULL default 'Diary entry',
date TIMESTAMP NOT NULL,
text TEXT NOT NULL,
user_id UUID NOT NULL,
CONSTRAINT fk_user_id
FOREIGN KEY (user_id)
REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE diary_entries;


