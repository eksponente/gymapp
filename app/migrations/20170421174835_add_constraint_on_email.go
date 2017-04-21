package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20170421174835, Down20170421174835)
}

// Up20170421174835 updates the database to the new requirements
func Up20170421174835(tx *sql.Tx) error {
	_, err := tx.Exec("ALTER TABLE users ADD CONSTRAINT unique_email UNIQUE (Email);")
	if err != nil {
		return err
	}
	return nil
}

// Down20170421174835 should send the database back to the state it was from before Up was ran
func Down20170421174835(tx *sql.Tx) error {
	return nil
}
