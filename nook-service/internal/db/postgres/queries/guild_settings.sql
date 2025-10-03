-- name: UpsertGuildSettings :exec
INSERT INTO guild_settings (
    guild_id,
    command_prefix,
    color_scheme,
    updated_at
) VALUES ($1, $2, $3, $4) 
ON CONFLICT (guild_id) DO UPDATE SET 
    command_prefix = EXCLUDED.command_prefix,
    color_scheme = EXCLUDED.color_scheme,
    updated_at = EXCLUDED.updated_at;

-- name: GetGuildSettings :one
SELECT * FROM guild_settings WHERE guild_id = $1;

-- name: DeleteGuildSettings :exec
DELETE FROM guild_settings WHERE guild_id = $1;

-- name: GetAllGuildSettings :many
SELECT * FROM guild_settings;
