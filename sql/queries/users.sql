-- name: CreateUser :one
INSERT INTO users (name, api_key, created_at, updated_at)
VALUES (
    ?,
    ?,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
RETURNING *;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUser :one
SELECT * FROM users
    WHERE name = ?;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: DeleteUser :exec
DELETE FROM users
    WHERE name = ?;

-- name: UpdateApiKey :exec
UPDATE users
    SET api_key = ?
    WHERE name = ?;