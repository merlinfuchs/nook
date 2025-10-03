-- name: UpsertGuild :one
INSERT INTO guilds (
    id, 
    name,
    description, 
    icon,
    unavailable,
    owner_user_id, 
    created_at, 
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (id) DO UPDATE SET 
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    icon = EXCLUDED.icon,
    unavailable = EXCLUDED.unavailable,
    owner_user_id = EXCLUDED.owner_user_id,
    updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: UpdateGuildDeleted :exec
UPDATE guilds SET deleted = $2 WHERE id = $1;

-- name: UpdateGuildUnavailable :exec
UPDATE guilds SET unavailable = $2 WHERE id = $1;

-- name: GetGuild :one
SELECT * FROM guilds WHERE id = $1;

-- name: GetGuildsByOwnerUserID :many
SELECT * FROM guilds WHERE owner_user_id = $1;
