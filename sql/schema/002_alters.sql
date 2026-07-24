-- +goose Up
CREATE TABLE alters(
id UUID PRIMARY KEY,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
avatar TEXT NOT NULL default 'assets/default_avater.jpg',
name TEXT NOT NULL,
pronouns TEXT NOT NULL default '',
age TEXT NOT NULL default '',
alter_role TEXT NOT NULL default '',
description TEXT NOT NULL default '',
colour TEXT NOT NULL default '000000',
fronting BOOL NOT NULL default false,
fronting_since TIMESTAMP default NULL,
user_id  UUID NOT NULL,
CONSTRAINT fk_user_id
FOREIGN KEY (user_id)
REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE alters;

