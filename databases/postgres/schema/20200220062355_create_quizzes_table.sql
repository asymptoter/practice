-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE SEQUENCE quizzes_id_seq;
CREATE TABLE IF NOT EXISTS quizzes (
    id INT NOT NULL DEFAULT nextval('quizzes_id_seq'),
    content VARCHAR(512),
    image_url VARCHAR(100),
    options VARCHAR(64) ARRAY,
    answer VARCHAR(64),
    creator UUID,
    category VARCHAR(64),
    PRIMARY KEY (id)
);
CREATE INDEX ON quizzes (creator, category);
ALTER SEQUENCE quizzes_id_seq OWNED BY quizzes.id;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS quizzes;
