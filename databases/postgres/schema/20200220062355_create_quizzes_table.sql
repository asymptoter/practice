-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS `quizzes` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `content` VARCHAR(512),
    `option1` CHAR(32),
    `option2` CHAR(32),
    `option3` CHAR(32),
    `option4` CHAR(32),
    `answer` INT(1),
    `creator` CHAR(36),
    PRIMARY KEY (`id`),
    KEY `creator` (`creator`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS quizzes;
