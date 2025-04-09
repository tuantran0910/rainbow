-- +goose Up
-- +goose StatementBegin
ALTER TABLE sellers DROP COLUMN user_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sellers ADD COLUMN user_id UUID NOT NULL;
-- +goose StatementEnd
