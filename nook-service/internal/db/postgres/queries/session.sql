-- name: CreateSession :exec
INSERT INTO sessions (
    key_hash, 
    user_id, 
    token_access,
    token_refresh, 
    token_scopes, 
    token_expires_at, 
    guilds,
    created_at,
    expires_at
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
);

-- name: GetSession :one
SELECT * FROM sessions WHERE key_hash = $1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE key_hash = $1;