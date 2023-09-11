CREATE DATABASE repositories_db;

\c repositories_db;

DROP TABLE IF EXISTS repositories;

CREATE TABLE repositories (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    owner VARCHAR(255) NOT NULL,
    creationDate DATE NOT NULL,
    configurationId VARCHAR(255),
    state VARCHAR(255) NOT NULL
);
