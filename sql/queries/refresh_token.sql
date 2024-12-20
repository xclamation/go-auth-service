-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash)
VALUES ($1, $2)
RETURNING id, user_id, token_hash, created_at;

-- name: GetRefreshTokenByHash :one
SELECT id, user_id, token_hash, created_at
FROM refresh_tokens
WHERE token_hash = $1;

-- name: GetRefreshTokenByUserID :many
SELECT token_hash
FROM refresh_tokens
WHERE user_id = $1;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE id = $1;