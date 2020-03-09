-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- CREATE EXTENSION IF NOT EXISTS 'uuid-ossp';
CREATE TABLE IF NOT EXISTS users (
    id UUID,
    token UUID,
    email VARCHAR(320) UNIQUE NOT NULL,
    password CHAR(60) NOT NULL, 
    register_date BIGINT,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX ON users (email);
CREATE UNIQUE INDEX ON users (token);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS users;
