-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE SEQUENCE games_id_seq;
CREATE TABLE IF NOT EXISTS games (
    id INT NOT NULL DEFAULT nextval('games_id_seq'),
    name VARCHAR(32),
    quiz_ids INT ARRAY,
    mode SMALLINT,
    count_down SMALLINT,
    creator UUID,
    PRIMARY KEY (id)
);
CREATE INDEX creator_name ON games (creator, name);
ALTER SEQUENCE games_id_seq OWNED BY games.id;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS games;
