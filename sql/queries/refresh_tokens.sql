-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $2,
    NOW(),
    NOW(),
    $1,
    NOW() + make_interval(0,0,0,60) 
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: RevokeToken :one
UPDATE refresh_tokens SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1
RETURNING *;