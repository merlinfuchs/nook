-- name: SetModuleValue :one
INSERT INTO module_values (
    guild_id,
    module_id,
    key,
    value,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) ON CONFLICT (guild_id, module_id, key) DO UPDATE SET
    value = EXCLUDED.value,
    updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: GetModuleValue :one
SELECT * FROM module_values WHERE guild_id = $1 AND module_id = $2 AND key = $3;

-- name: GetModuleValueForUpdate :one
SELECT * FROM module_values WHERE guild_id = $1 AND module_id = $2 AND key = $3 FOR UPDATE;

-- name: DeleteModuleValue :exec
DELETE FROM module_values WHERE guild_id = $1 AND module_id = $2 AND key = $3;
