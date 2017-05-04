package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20170421160454, Down20170421160454)
}

// Up20170421160454 updates the database to the new requirements
func Up20170421160454(tx *sql.Tx) error {
	_, err := tx.Exec("ALTER TABLE users ADD CONSTRAINT unique_usernames UNIQUE (Username);")
	if err != nil {
		return err
	}
	return nil
}

// Down20170421160454 should send the database back to the state it was from before Up was ran
func Down20170421160454(tx *sql.Tx) error {
	_, err := tx.Exec("ALTER TABLE users DROP CONSTRAINT IF EXISTS unique_usernames;")
	if err != nil {
		return err
	}
	return nil
}
