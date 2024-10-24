-- name: CreateChirp :one
INSERT INTO chirps(id, created_at, updated_at, body, user_id)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *;

-- name: DeleteChirps :exec
DELETE FROM chirps;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;

-- name: GetChirps_ASC :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetChirps_DESC :many
SELECT * FROM chirps
ORDER BY created_at DESC;

-- name: GetChirpsByAuthorID_ASC :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: GetChirpsByAuthorID_DESC :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetChirp :one
SELECT * FROM chirps
WHERE id = $1;