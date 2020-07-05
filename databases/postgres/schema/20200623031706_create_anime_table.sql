-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS animes (
    id UUID,
    name VARCHAR NOT NULL,
    author VARCHAR NOT NULL,
    director VARCHAR NOT NULL,
    produce VARCHAR NOT NULL,
    PRIMARY KEY (id)
)
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
