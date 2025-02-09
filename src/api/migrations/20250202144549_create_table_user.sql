-- +goose Up
-- +goose StatementBegin
CREATE TYPE user_role AS ENUM ('USER', 'ADMIN');

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL CHECK (LENGTH(password) > 6 AND LENGTH(password) <= 255),
    first_name VARCHAR(255) NOT NULL CHECK (LENGTH(first_name) > 0 AND LENGTH(first_name) <= 255),
    last_name VARCHAR(255) NOT NULL CHECK (LENGTH(last_name) > 0 AND LENGTH(last_name) <= 255),
    last_login TIMESTAMP NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    phone_number VARCHAR(255) NOT NULL CHECK (LENGTH(phone_number) = 10),
    role user_role NOT NULL DEFAULT 'USER',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS user_role;
-- +goose StatementEnd
