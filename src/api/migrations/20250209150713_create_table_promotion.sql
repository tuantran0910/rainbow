-- +goose Up
-- +goose StatementBegin
CREATE TYPE discount_type AS ENUM ('PERCENTAGE', 'FIXED');

CREATE TABLE IF NOT EXISTS promotions (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    discount_type discount_type NOT NULL,
    discount_value DECIMAL(10, 2) NOT NULL CHECK (discount_value >= 0),
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL CHECK (end_date > start_date),
    max_uses INT NOT NULL CHECK (max_uses >= 0),
    used_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS promotions;

DROP TYPE IF EXISTS discount_type;
-- +goose StatementEnd
