-- name: UpsertModuleSettings :exec
INSERT INTO module_settings (
    guild_id,
    module_id,
    enabled,
    command_overwrites,
    config,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6) 
ON CONFLICT (guild_id, module_id) DO UPDATE SET 
    enabled = EXCLUDED.enabled,
    command_overwrites = EXCLUDED.command_overwrites,
    config = EXCLUDED.config,
    updated_at = EXCLUDED.updated_at;

-- name: GetModuleSettings :one
SELECT * FROM module_settings WHERE guild_id = $1 AND module_id = $2;

-- name: DeleteModuleSettings :exec
DELETE FROM module_settings WHERE guild_id = $1 AND module_id = $2;

-- name: GetModuleSettingsByGuildID :many
SELECT * FROM module_settings WHERE guild_id = $1;
