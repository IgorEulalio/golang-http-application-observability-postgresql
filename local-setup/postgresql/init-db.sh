CREATE TABLE IF NOT EXISTS repositories (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    owner VARCHAR(255) NOT NULL,
    creationDate DATE NOT NULL,
    configurationId VARCHAR(255)
);