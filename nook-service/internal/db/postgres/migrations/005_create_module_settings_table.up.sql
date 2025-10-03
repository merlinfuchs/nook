CREATE TABLE IF NOT EXISTS module_settings (
    guild_id BIGINT REFERENCES guilds(id) ON DELETE CASCADE,
    module_id TEXT,

    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    command_overwrites JSONB NOT NULL DEFAULT '{}',
    config JSONB NOT NULL,

    updated_at TIMESTAMP NOT NULL,

    PRIMARY KEY (guild_id, module_id)
);
