package models

import "time"

type Token struct {
	UserID         int       `db:"user_id"`
	Token          string    `db:"token"`
	ExpirationDate time.Time `db:"expirationdate"`
}
