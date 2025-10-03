CREATE TABLE IF NOT EXISTS module_values (
    id BIGSERIAL PRIMARY KEY,
    guild_id BIGINT NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
    module_id TEXT NOT NULL,

    key TEXT NOT NULL,
    value JSONB NOT NULL,
   
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    UNIQUE (guild_id, module_id, key)
);

CREATE INDEX IF NOT EXISTS module_values_guild_id ON module_values (guild_id);
CREATE INDEX IF NOT EXISTS module_values_module_id ON module_values (module_id);
CREATE INDEX IF NOT EXISTS module_values_module_key ON module_values (key);
