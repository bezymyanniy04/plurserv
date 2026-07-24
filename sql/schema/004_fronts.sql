-- +goose Up
CREATE TABLE fronts(
id UUID PRIMARY KEY,
started_at TIMESTAMP NOT NULL,
ended_at TIMESTAMP default NULL,
alter_id  UUID NOT NULL,
CONSTRAINT fk_alter_id
FOREIGN KEY (alter_id)
REFERENCES alters(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE fronts;
