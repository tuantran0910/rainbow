-- +goose Up
-- +goose StatementBegin
INSERT INTO users (
    id,
    email,
    password,
    first_name,
    last_name,
    phone_number,
    role,
    is_active,
    last_login,
    created_at,
    updated_at
) VALUES (
    GEN_RANDOM_UUID(),
    'admin@example.com',
    -- The password is 'Admin123!' - this will be hashed properly
    '$2b$12$JWWHao5ZdvnuyfjARp7tBOeV3vNy3Y1zAQZiYVwPn19F4Zq06Zdwm',
    'Admin',
    'User',
    '1234567890',
    'ADMIN',
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT (email) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE email = 'admin@example.com';
-- +goose StatementEnd
