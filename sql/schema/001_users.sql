-- +goose Up
CREATE TABLE users(
    id UUID primary key,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT UNIQUE NOT NULL,
    hashed_password  TEXT NOT NULL DEFAULT 'unset',
    is_chirpy_red BOOLEAN NOT NULL DEFAULT false
);

-- +goose Down
DROP TABLE users;