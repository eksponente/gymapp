package migration

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/pressly/goose"
	"github.com/revel/revel"
)

func init() {
	goose.AddMigration(Up20170420220207, Down20170420220207)
}

// Up20170420220207 updates the database to the new requirements
func Up20170420220207(tx *sql.Tx) error {
	password, found := revel.Config.String("superuserPassword")
	fmt.Println(password, found)
	if !found {
		panic("Superuser password not set in config.")
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err := tx.Exec("INSERT INTO users (Name, Username, Password) VALUES ($1, $2, $3);", "Admin", "Admin", hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

// Down20170420220207 should send the database back to the state it was from before Up was ran
func Down20170420220207(tx *sql.Tx) error {
	return nil
}
