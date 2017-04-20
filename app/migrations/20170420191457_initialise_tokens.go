package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20170420191457, Down20170420191457)
}

// Up20170420191457 updates the database to the new requirements
func Up20170420191457(tx *sql.Tx) error {
	_, err := tx.Exec("CREATE TABLE tokens (user_id int NOT NULL, Token text NOT NULL, ExpirationDate timestamp with time zone NOT NULL, PRIMARY KEY(Token));")
	if err != nil {
		return err
	}
	return nil
}

// Down20170420191457 should send the database back to the state it was from before Up was ran
func Down20170420191457(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE tokens;")
	if err != nil {
		return err
	}
	return nil
}
