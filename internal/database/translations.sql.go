// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: translations.sql

package database

import (
	"context"
)

const createTranslation = `-- name: CreateTranslation :one
INSERT INTO translations (name, api_id, language_id, user_id)
VALUES(
    ?,
    ?,
    ?,
    ?
)
RETURNING id, name, api_id, language_id, user_id
`

type CreateTranslationParams struct {
	Name       string
	ApiID      string
	LanguageID string
	UserID     int64
}

func (q *Queries) CreateTranslation(ctx context.Context, arg CreateTranslationParams) (Translation, error) {
	row := q.db.QueryRowContext(ctx, createTranslation,
		arg.Name,
		arg.ApiID,
		arg.LanguageID,
		arg.UserID,
	)
	var i Translation
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ApiID,
		&i.LanguageID,
		&i.UserID,
	)
	return i, err
}

const deleteTranslationForUser = `-- name: DeleteTranslationForUser :exec
DELETE FROM translations
    WHERE user_id = ? AND api_id = ?
`

type DeleteTranslationForUserParams struct {
	UserID int64
	ApiID  string
}

func (q *Queries) DeleteTranslationForUser(ctx context.Context, arg DeleteTranslationForUserParams) error {
	_, err := q.db.ExecContext(ctx, deleteTranslationForUser, arg.UserID, arg.ApiID)
	return err
}

const getTranslationForUserById = `-- name: GetTranslationForUserById :one
SELECT id, name, api_id, language_id, user_id FROM translations
    WHERE api_id = ?
    AND user_id = ?
`

type GetTranslationForUserByIdParams struct {
	ApiID  string
	UserID int64
}

func (q *Queries) GetTranslationForUserById(ctx context.Context, arg GetTranslationForUserByIdParams) (Translation, error) {
	row := q.db.QueryRowContext(ctx, getTranslationForUserById, arg.ApiID, arg.UserID)
	var i Translation
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ApiID,
		&i.LanguageID,
		&i.UserID,
	)
	return i, err
}

const getTranslationsForUser = `-- name: GetTranslationsForUser :many
SELECT id, name, api_id, language_id, user_id FROM translations
    WHERE user_id = ?
`

func (q *Queries) GetTranslationsForUser(ctx context.Context, userID int64) ([]Translation, error) {
	rows, err := q.db.QueryContext(ctx, getTranslationsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Translation
	for rows.Next() {
		var i Translation
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ApiID,
			&i.LanguageID,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
