CREATE TABLE IF NOT EXISTS routes (
    id SERIAL PRIMARY KEY,
    path varchar(2048) NOT NULL UNIQUE,
    description TEXT NULL
);

