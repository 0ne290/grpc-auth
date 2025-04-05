CREATE TABLE users (
    uuid UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    name TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);