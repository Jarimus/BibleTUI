// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"time"
)

type Translation struct {
	ID         int64
	Name       string
	ApiID      string
	LanguageID string
	UserID     int64
}

type User struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
