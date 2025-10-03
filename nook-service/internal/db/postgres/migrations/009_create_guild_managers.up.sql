CREATE TABLE IF NOT EXISTS guild_managers (
    guild_id BIGINT NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    PRIMARY KEY (guild_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_guild_managers_guild_id ON guild_managers (guild_id);
CREATE INDEX IF NOT EXISTS idx_guild_managers_user_id ON guild_managers (user_id);
