-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sellers (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL CHECK (LENGTH(name) > 0 AND LENGTH(name) <= 255),
    link VARCHAR(255) NOT NULL,
    logo VARCHAR(255) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sellers;
-- +goose StatementEnd
