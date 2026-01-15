-- +goose Up

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    uuid UUID PRIMARY KEY,
    login VARCHAR(32) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS games(
    uuid UUID PRIMARY KEY,
    field JSONB NOT NULL,
    status INTEGER NOT NULL DEFAULT 0,
    player_x UUID NOT NULL,
    player_o UUID,
    current_turn UUID NOT NULL,
    symbols JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS refresh_tokens(
    token_hash TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_tokens;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS games;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd