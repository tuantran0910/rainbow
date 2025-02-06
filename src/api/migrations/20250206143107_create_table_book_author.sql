-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS book_authors (
    id UUID PRIMARY KEY,
    book_id UUID NULL REFERENCES books (id) ON DELETE CASCADE,
    author_id UUID NULL REFERENCES authors (id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create indexes for foreign keys
CREATE INDEX IF NOT EXISTS idx_book_authors_book_id ON book_authors (book_id);
CREATE INDEX IF NOT EXISTS idx_book_authors_author_id ON book_authors (author_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS book_authors;
-- +goose StatementEnd
