-- +goose Up
-- +goose StatementBegin
ALTER TABLE books
ADD COLUMN sold_count INT CHECK (sold_count >= 0) DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE books
DROP COLUMN sold_count;
-- +goose StatementEnd
