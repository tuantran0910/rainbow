-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_promotions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    promotion_id UUID NOT NULL REFERENCES promotions (id) ON DELETE CASCADE,
    used_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT user_promotion_unique UNIQUE (user_id, promotion_id)
);

CREATE INDEX user_promotions_user_id_idx ON user_promotions (user_id);
CREATE INDEX user_promotions_promotion_id_idx ON user_promotions (promotion_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_promotions;
-- +goose StatementEnd
