-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: UpsertUser :one
INSERT INTO users (
    id,
    username,
    discriminator,
    display_name,
    avatar,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) ON CONFLICT (id) 
DO UPDATE SET 
    username = EXCLUDED.username,
    discriminator = EXCLUDED.discriminator,
    display_name = EXCLUDED.display_name,
    avatar = EXCLUDED.avatar,
    updated_at = EXCLUDED.updated_at
RETURNING *;
