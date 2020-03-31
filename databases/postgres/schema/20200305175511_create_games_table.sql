-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS games (
    id UUID,
    name VARCHAR(32),
    quiz_ids INT ARRAY,
    mode SMALLINT,
    count_down SMALLINT,
    creator UUID,
    PRIMARY KEY (id)
);
CREATE INDEX creator_name_idx ON games (creator, name);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS games;
