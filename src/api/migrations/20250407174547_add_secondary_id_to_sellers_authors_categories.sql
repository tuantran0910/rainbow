-- +goose Up
-- +goose StatementBegin
ALTER TABLE sellers
ADD COLUMN secondary_id VARCHAR(255) UNIQUE DEFAULT NULL;

ALTER TABLE authors
ADD COLUMN secondary_id VARCHAR(255) UNIQUE DEFAULT NULL;

ALTER TABLE categories
ADD COLUMN secondary_id VARCHAR(255) UNIQUE DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sellers
DROP COLUMN secondary_id;

ALTER TABLE authors
DROP COLUMN secondary_id;

ALTER TABLE categories
DROP COLUMN secondary_id;
-- +goose StatementEnd
