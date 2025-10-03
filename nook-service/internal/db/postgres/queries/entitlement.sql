-- name: UpsertSubscriptionEntitlement :one
INSERT INTO entitlements (
    id,
    guild_id,
    subscription_id,
    type,
    plan_ids,
    starts_at,
    ends_at,
    created_at,
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
    $9
) ON CONFLICT (guild_id, subscription_id) DO UPDATE SET
    type = EXCLUDED.type,
    plan_ids = EXCLUDED.plan_ids,
    starts_at = EXCLUDED.starts_at,
    ends_at = EXCLUDED.ends_at,
    updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: GetEntitlementsBySubscriptionID :many
SELECT * FROM entitlements WHERE subscription_id = $1;

-- name: GetEntitlementsByGuildID :many
SELECT * FROM entitlements WHERE guild_id = $1;

-- name: GetActiveEntitlementsByGuildID :many
SELECT * FROM entitlements 
WHERE guild_id = $1 
AND (starts_at IS NULL OR starts_at <= $2) 
AND (ends_at IS NULL OR ends_at > $2);
