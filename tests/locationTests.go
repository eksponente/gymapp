package tests

import (
	"errors"
	"gymapp/app/controllers"

	"github.com/pressly/goose"
	db "github.com/revel/modules/db/app"
	"github.com/revel/revel"
	"github.com/revel/revel/testing"
)

type LocationTest struct {
	testing.TestSuite
}

func (t *LocationTest) Before() {
	println("Set up")
	if name, _ := revel.Config.String("db.name"); name != "gymapptest" {
		panic(errors.New("Not connected to test database.  RDS_DB_NAME must be gymapptest."))
	}
	println("MIGRATING DOWN")
	goose.DownTo(db.Db, "../app/migrations", 20170419150037)
	goose.Down(db.Db, "../app/migrations")
	println("MIGRATING UP")
	goose.Up(db.Db, "../app/migrations")

	controllers.Dbm.SelectOne(&user, "SELECT * FROM users WHERE Email=$1", "rugilena@gmail.com")
}

func (t *LocationTest) After() {
	println("Tear down")
}

func (t *LocationTest) TestCreatingNewLocation() {
	txn, _ := controllers.Dbm.Begin()
	rows, err := controllers.CreateLocation(user.UserId, -64.3333, 18.3333345, "Dummy address", "Dummy name", txn)
	if err != nil {
		panic(err)
	}
	t.Assert(rows == 1)
}
