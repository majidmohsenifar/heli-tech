CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    code varchar(128) NOT NULL UNIQUE
);

