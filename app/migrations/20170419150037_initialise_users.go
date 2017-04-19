package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20170419150037, Down20170419150037)
}

// Up20170419150037 updates the database to the new requirements
func Up20170419150037(tx *sql.Tx) error {
	_, err := tx.Exec("CREATE TABLE users (user_id int NOT NULL, Name text, Username text,Password text, PRIMARY KEY(user_id));")
	if err != nil {
		return err
	}
	return nil
}

// Down20170419150037 should send the database back to the state it was from before Up was ran
func Down20170419150037(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE users;")
	if err != nil {
		return err
	}
	return nil
}
