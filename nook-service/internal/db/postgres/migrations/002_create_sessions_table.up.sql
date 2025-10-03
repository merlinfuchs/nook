CREATE TABLE IF NOT EXISTS sessions (
    key_hash TEXT PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_access TEXT NOT NULL,
    token_refresh TEXT NOT NULL,
    token_scopes TEXT[] NOT NULL,
    token_expires_at TIMESTAMP NOT NULL,
    guilds JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL
);
