-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_sellers_secondary_id ON sellers (secondary_id);
CREATE INDEX IF NOT EXISTS idx_authors_secondary_id ON authors (secondary_id);
CREATE INDEX IF NOT EXISTS idx_categories_secondary_id ON categories (secondary_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_sellers_secondary_id;
DROP INDEX IF EXISTS idx_authors_secondary_id;
DROP INDEX IF EXISTS idx_categories_secondary_id;
-- +goose StatementEnd
