-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE users ADD activated BIT(1) NOT NULL;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE users DROP COLUMN activated;
