// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package sqlc

import (
	"database/sql"
)

type Product struct {
	ID     string
	Name   string
	Code   string
	Price  int64
	UserID string
}

type User struct {
	ID       string
	Name     string
	Email    string
	Bio      sql.NullString
	Password string
}