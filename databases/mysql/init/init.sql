CREATE DATABASE IF NOT EXISTS practice;
USE practice;

CREATE TABLE IF NOT EXISTS `users` (
    `id` CHAR(36), 
    `email` VARCHAR(320),
    `password` CHAR(60) BINARY,
    `token` CHAR(36),
    PRIMARY KEY (`id`),
    UNIQUE KEY `email` (`email`),
    UNIQUE KEY `token` (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
