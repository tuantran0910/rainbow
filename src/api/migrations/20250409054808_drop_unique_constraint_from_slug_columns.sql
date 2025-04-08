-- +goose Up
-- +goose StatementBegin
-- Drop unique constraint from slug column in authors table
ALTER TABLE authors DROP CONSTRAINT IF EXISTS authors_slug_key;

-- Drop unique constraint from slug column in categories table
ALTER TABLE categories DROP CONSTRAINT IF EXISTS categories_slug_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Add unique constraint back to slug column in authors table
ALTER TABLE authors ADD CONSTRAINT authors_slug_key UNIQUE (slug);

-- Add unique constraint back to slug column in categories table
ALTER TABLE categories ADD CONSTRAINT categories_slug_key UNIQUE (slug);
-- +goose StatementEnd
