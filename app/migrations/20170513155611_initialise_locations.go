package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20170513155611, Down20170513155611)
}

// Up20170513155611 updates the database to the new requirements
func Up20170513155611(tx *sql.Tx) error {
	_, err := tx.Exec("CREATE TABLE locations(id serial primary key, coordinates geometry(POINT,4326) NOT NULL, address text NOT NULL, location_name text NOT NULL, user_id int NOT NULL, lat double precision NOT NULL, lon double precision NOT NULL, FOREIGN KEY(user_id) REFERENCES users(user_id));")
	if err != nil {
		return err
	}
	return nil
}

// Down20170513155611 should send the database back to the state it was from before Up was ran
func Down20170513155611(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE locations;")
	if err != nil {
		return err
	}
	return nil
}
