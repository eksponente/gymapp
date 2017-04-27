package models

type User struct {
	UserId      int `db:"user_id"`
	Name        string
	Password    string
	Email       string
	IsSuperuser bool
}
