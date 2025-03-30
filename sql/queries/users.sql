-- name: CreateUser :one
INSERT INTO users (name, created_at, updated_at)
VALUES (
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