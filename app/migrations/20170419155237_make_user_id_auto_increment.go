package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20170419155237, Down20170419155237)
}

// Up20170419155237 updates the database to the new requirements
func Up20170419155237(tx *sql.Tx) error {
	_, err := tx.Exec("CREATE SEQUENCE IF NOT EXISTS user_id_seq; ALTER TABLE users ALTER user_id SET DEFAULT nextval('user_id_seq'); SELECT setval('user_id_seq', 1);")
	if err != nil {
		return err
	}
	return nil
}

// Down20170419155237 should send the database back to the state it was from before Up was ran
func Down20170419155237(tx *sql.Tx) error {
	return nil
}
