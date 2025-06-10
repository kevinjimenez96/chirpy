-- name: CreateChirp :one

INSERT INTO chirps (id, created_at, updated_at, user_id, body)
VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2) RETURNING *;

-- name: GetAllChirps :many

SELECT *
FROM chirps
ORDER BY
    CASE WHEN @sort::text = 'ASC' THEN created_at END ASC,
    CASE WHEN @sort::text = 'DESC' THEN created_at END DESC;

-- name: GetAllChirpsByAuthor :many

SELECT *
FROM chirps
WHERE user_id = $1
ORDER BY
    CASE WHEN @sort::text = 'ASC' THEN created_at END ASC,
    CASE WHEN @sort::text = 'DESC' THEN created_at END DESC;

-- name: GetChirpById :one

SELECT *
FROM chirps
WHERE id = $1;

-- name: DeleteChirpById :one
DELETE
FROM chirps
WHERE id = $1 AND user_id = $2
RETURNING id;