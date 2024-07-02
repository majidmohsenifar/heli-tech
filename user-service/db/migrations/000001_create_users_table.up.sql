CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email varchar(320) NOT NULL UNIQUE,
    password varchar(128) NOT NULL
);
