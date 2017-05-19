package models

type Location struct {
	UserID    int    `db:"user_id"`
	Latitude  string `db:"lat"`
	Longitude string `db:"lon"`
	ID        int    `db:"id"`
	Address   string `db:"address"`
	Name      string `db:"name"`
}
