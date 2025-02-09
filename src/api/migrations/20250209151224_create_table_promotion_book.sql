-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS promotion_books (
    id UUID PRIMARY KEY,
    promotion_id UUID NOT NULL REFERENCES promotions (id) ON DELETE SET NULL,
    book_id UUID NOT NULL REFERENCES books (id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX promotion_books_promotion_id_idx ON promotion_books (promotion_id);
CREATE INDEX promotion_books_book_id_idx ON promotion_books (book_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS promotion_books;
-- +goose StatementEnd
