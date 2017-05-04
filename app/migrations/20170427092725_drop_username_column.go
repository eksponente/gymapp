package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20170427092725, Down20170427092725)
}

// Up20170427092725 updates the database to the new requirements
func Up20170427092725(tx *sql.Tx) error {
	_, err := tx.Exec("ALTER TABLE users DROP COLUMN username;")
	if err != nil {
		return err
	}
	return nil
}

// Down20170427092725 should send the database back to the state it was from before Up was ran
func Down20170427092725(tx *sql.Tx) error {
	_, err := tx.Exec("ALTER TABLE users ADD COLUMN username text;")
	if err != nil {
		return err
	}
	return nil
}
