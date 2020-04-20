-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS game_results (
    user_id UUID,
    game_id UUID,
    play_date BIGINT,
    correct_count INT,
    time_spent BIGINT,
    PRIMARY KEY (user_id, game_id, play_date)
);
CREATE INDEX ON game_results (game_id);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS game_results;
