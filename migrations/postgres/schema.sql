CREATE TABLE IF NOT EXISTS users (
    uuid UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    name TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
    refresh_token UUID PRIMARY KEY,
    user_uuid UUID REFERENCES users(uuid) ON DELETE CASCADE NOT NULL,
    expiration_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks (
    uuid UUID PRIMARY KEY,
    user_uuid UUID REFERENCES users(uuid) ON DELETE CASCADE NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);