-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS books (
    id UUID PRIMARY KEY,
    secondary_id VARCHAR(255) UNIQUE DEFAULT NULL,
    category_id UUID NOT NULL REFERENCES categories (id) ON DELETE SET NULL,
    seller_id UUID NOT NULL REFERENCES sellers (id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL CHECK (LENGTH(name) > 0 AND LENGTH(name) <= 255),
    description TEXT DEFAULT NULL,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    original_price DECIMAL(10, 2) CHECK (original_price >= 0 AND original_price >= price),
    rating_average DECIMAL(2, 1) CHECK (rating_average BETWEEN 0 AND 5),
    review_count INT CHECK (review_count >= 0),
    page_count INT CHECK (page_count >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS inventories (
    id UUID PRIMARY KEY,
    book_id UUID NOT NULL REFERENCES books (id) ON DELETE CASCADE,
    stock INT NOT NULL CHECK (stock >= 0),
    last_restocked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create indexes for foreign keys
CREATE INDEX IF NOT EXISTS idx_inventories_book_id ON inventories (book_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS inventories;
DROP TABLE IF EXISTS books;
-- +goose StatementEnd
