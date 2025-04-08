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
    '$2a$10$UUE7wP9Cz8qVYYZfX2OXxOdwLYS4UOZPzUSo0UHsQqopcW7UYIcpC',
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
