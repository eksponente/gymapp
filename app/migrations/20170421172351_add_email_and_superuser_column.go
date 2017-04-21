package migration

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/pressly/goose"
	"github.com/revel/revel"
)

func init() {
	goose.AddMigration(Up20170421172351, Down20170421172351)
}

// Up20170421172351 updates the database to the new requirements
func Up20170421172351(tx *sql.Tx) error {
	password, found := revel.Config.String("superuserPassword")
	fmt.Println(password, found)
	if !found {
		panic("Superuser password not set in config.")
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err := tx.Exec("ALTER TABLE users ADD COLUMN Email text;")
	if err != nil {
		return err
	}
	_, err = tx.Exec("ALTER TABLE users ADD COLUMN IsSuperuser boolean DEFAULT FALSE;")
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO users (Name, Username, Email, IsSuperuser, Password) VALUES ($1, $2, $3, $4, $5);", "Rugile", "eksponente", "rugilena@gmail.com", true, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

// Down20170421172351 should send the database back to the state it was from before Up was ran
func Down20170421172351(tx *sql.Tx) error {
	return nil
}
