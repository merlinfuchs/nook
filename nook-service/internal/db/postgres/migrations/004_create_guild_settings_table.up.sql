CREATE TABLE IF NOT EXISTS guild_settings (
    guild_id BIGINT PRIMARY KEY REFERENCES guilds(id) ON DELETE CASCADE,

    command_prefix TEXT,
    color_scheme TEXT,

    updated_at TIMESTAMP NOT NULL
);
