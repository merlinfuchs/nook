CREATE TABLE IF NOT EXISTS subscriptions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    status TEXT NOT NULL,
    paddle_subscription_id TEXT NOT NULL UNIQUE,
    paddle_customer_id TEXT NOT NULL,
    paddle_product_ids TEXT[] NOT NULL,
    paddle_price_ids TEXT[] NOT NULL,

    created_at TIMESTAMP NOT NULL,
    started_at TIMESTAMP NOT NULL,
    paused_at TIMESTAMP,
    canceled_at TIMESTAMP,
    current_period_ends_at TIMESTAMP,
    current_period_starts_at TIMESTAMP,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions (user_id);
