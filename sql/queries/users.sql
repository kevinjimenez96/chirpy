-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users SET email = $2, hashed_password= $3, updated_at = NOW()
WHERE id = $1
RETURNING *;


-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: UpdateUserIsChirpyRedById :one
UPDATE users SET is_chirpy_red = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;