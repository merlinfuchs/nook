-- name: UpsertGuildManager :one
INSERT INTO guild_managers (
    guild_id,
    user_id,
    role,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (guild_id, user_id) DO UPDATE SET
    role = EXCLUDED.role,
    updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: GetGuildManager :one
SELECT * FROM guild_managers WHERE guild_id = $1 AND user_id = $2;

-- name: DeleteGuildManager :exec
DELETE FROM guild_managers WHERE guild_id = $1 AND user_id = $2;

-- name: GetGuildManagers :many
SELECT * FROM guild_managers WHERE guild_id = $1;

-- name: GetGuildManagersWithUsers :many
SELECT sqlc.embed(guild_managers), sqlc.embed(users) 
FROM guild_managers JOIN users ON guild_managers.user_id = users.id 
WHERE guild_managers.guild_id = $1;
