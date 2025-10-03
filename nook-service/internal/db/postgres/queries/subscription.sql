-- name: UpsertSubscription :one
INSERT INTO subscriptions (
    id,
    user_id,
    status,
    paddle_subscription_id,
    paddle_customer_id,
    paddle_product_ids,
    paddle_price_ids,
    created_at,
    started_at,
    paused_at,
    canceled_at,
    current_period_ends_at,
    current_period_starts_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13,
    $14
) ON CONFLICT (paddle_subscription_id) DO UPDATE SET
    status = EXCLUDED.status,
    paddle_customer_id = EXCLUDED.paddle_customer_id,
    paddle_product_ids = EXCLUDED.paddle_product_ids,
    paddle_price_ids = EXCLUDED.paddle_price_ids,
    started_at = EXCLUDED.started_at,
    paused_at = EXCLUDED.paused_at,
    canceled_at = EXCLUDED.canceled_at,
    current_period_ends_at = EXCLUDED.current_period_ends_at,
    current_period_starts_at = EXCLUDED.current_period_starts_at,
    updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: GetSubscriptionsByUserID :many
SELECT * FROM subscriptions WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetSubscriptionsByGuildID :many
SELECT subscriptions.* FROM subscriptions 
LEFT JOIN entitlements ON subscriptions.id = entitlements.subscription_id 
WHERE entitlements.guild_id = $1
ORDER BY subscriptions.created_at DESC;

-- name: GetSubscriptionByID :one
SELECT * FROM subscriptions WHERE id = $1;
