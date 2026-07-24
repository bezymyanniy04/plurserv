-- +goose Up
CREATE TABLE users(
id UUID PRIMARY KEY,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
email TEXT UNIQUE NOT NULL,
hashed_password TEXT NOT NULL,
avatar TEXT NOT NULL default 'assets/default_avater.jpg',
system_name TEXT NOT NULL default 'System',
theme INTEGER NOT NULL DEFAULT 0,
font TEXT NOT NULL default 'Arial' 
);

-- +goose Down
DROP TABLE users;
