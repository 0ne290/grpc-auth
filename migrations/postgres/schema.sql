CREATE TABLE users (
    uuid UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    name TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE sessions (
    refresh_token UUID PRIMARY KEY,
    user_uuid UUID REFERENCES users(uuid) ON DELETE CASCADE NOT NULL,
    expiration_at TIMESTAMP NOT NULL
);