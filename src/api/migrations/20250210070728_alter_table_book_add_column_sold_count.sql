-- +goose Up
-- +goose StatementBegin
ALTER TABLE books
ADD COLUMN sold_count INT NOT NULL DEFAULT 0 CHECK (sold_count >= 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE books
DROP COLUMN sold_count;
-- +goose StatementEnd
