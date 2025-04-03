-- name: CreateTranslation :one
INSERT INTO translations (name, api_id, language_id, user_id)
VALUES(
    ?,
    ?,
    ?,
    ?
)
RETURNING *;

-- name: GetTranslationsForUser :many
SELECT * FROM translations
    WHERE user_id = ?;

-- name: DeleteTranslationForUser :exec
DELETE FROM translations
    WHERE user_id = ? AND api_id = ?;

-- name: GetTranslationForUserById :one
SELECT * FROM translations
    WHERE api_id = ?
    AND user_id = ?;