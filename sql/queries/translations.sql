-- name: CreateTranslation :one
INSERT INTO translations (name, user_id, api_id)
VALUES(
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