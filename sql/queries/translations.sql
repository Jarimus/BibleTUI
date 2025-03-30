-- name: CreateTranslation :one
INSERT INTO translations (name, api_id, user_id)
VALUES(
    ?,
    ?,
    ?
)
RETURNING *;

-- name: GetTranslationsForUser :many
SELECT * FROM translations
    WHERE user_id = ?;