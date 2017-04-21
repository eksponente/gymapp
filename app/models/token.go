package models

type Token struct {
	UserID         int    `db:"user_id"`
	Token          string `db:"token"`
	ExpirationDate string
}
