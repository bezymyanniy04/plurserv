-- +goose Up
CREATE TABLE diaries(
id UUID PRIMARY KEY,
alter_id  UUID NOT NULL UNIQUE,
CONSTRAINT fk_alter_id
FOREIGN KEY (alter_id)
REFERENCES alters(id) ON DELETE CASCADE,
bg_colour TEXT NOT NULL default '255 255 255',
bg_colour2 TEXT NOT NULL default '255 255 255',
block_colour TEXT NOT NULL default '120 120 120',
text_colour TEXT NOT NULL default '0 0 0',
font TEXT NOT NULL default 'Arial',
user_id UUID NOT NULL,
CONSTRAINT fk_user_id FOREIGN KEY (user_id) 
REFERENCES users(id) ON DELETE CASCADE,
name TEXT NOT NULL
);

-- +goose Down
DROP TABLE diaries;


