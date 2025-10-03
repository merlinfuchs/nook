CREATE TABLE IF NOT EXISTS entitlements (
    id BIGINT PRIMARY KEY,
    guild_id BIGINT NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
    subscription_id BIGINT REFERENCES subscriptions(id) ON DELETE CASCADE,

    type TEXT NOT NULL, -- "subscription", "manual"
    plan_ids TEXT[] NOT NULL,

    starts_at TIMESTAMP,
    ends_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    UNIQUE (guild_id, subscription_id)
);

CREATE INDEX IF NOT EXISTS entitlements_subscription_id ON entitlements (subscription_id);
CREATE INDEX IF NOT EXISTS entitlements_guild_id ON entitlements (guild_id);
