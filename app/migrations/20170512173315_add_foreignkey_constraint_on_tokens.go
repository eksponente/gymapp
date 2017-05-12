package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20170512173315, Down20170512173315)
}

// Up20170512173315 updates the database to the new requirements
func Up20170512173315(tx *sql.Tx) error {
	_, err := tx.Exec("ALTER TABLE tokens ADD FOREIGN KEY(user_id) REFERENCES users(user_id);")
	if err != nil {
		return err
	}
	return nil
}

// Down20170512173315 should send the database back to the state it was from before Up was ran
func Down20170512173315(tx *sql.Tx) error {
	_, err := tx.Exec("alter table tokens drop constraint tokens_user_id_fkey;")
	if err != nil {
		return err
	}
	return nil
}
