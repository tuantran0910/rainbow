-- +goose Up
-- +goose StatementBegin
-- Drop check constraint from slug column in authors table
ALTER TABLE authors DROP CONSTRAINT IF EXISTS authors_slug_check;

-- Drop check constraint from slug column in categories table
ALTER TABLE categories DROP CONSTRAINT IF EXISTS categories_slug_check;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Add check constraint back to slug column in authors table
ALTER TABLE authors ADD CONSTRAINT authors_slug_check CHECK (LENGTH(slug) > 0 AND LENGTH(slug) <= 255);

-- Add check constraint back to slug column in categories table
ALTER TABLE categories ADD CONSTRAINT categories_slug_check CHECK (LENGTH(slug) > 0 AND LENGTH(slug) <= 255);
-- +goose StatementEnd
