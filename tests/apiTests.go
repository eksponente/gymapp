package tests

import (
	"errors"
	"fmt"
	"gymapp/app/controllers"
	"gymapp/app/models"
	"net/url"

	"github.com/pressly/goose"
	db "github.com/revel/modules/db/app"
	"github.com/revel/revel"
	"github.com/revel/revel/testing"
)

type ApiTest struct {
	testing.TestSuite
}

func (t *ApiTest) Before() {
	println("Set up")
	if name, _ := revel.Config.String("db.name"); name != "gymapptest" {
		panic(errors.New("Not connected to test database. RDS_DB_NAME must be gymapptest."))
	}
	println("MIGRATING DOWN")
	goose.DownTo(db.Db, "../app/migrations", 20170419150037)
	goose.Down(db.Db, "../app/migrations")
	println("MIGRATING UP")
	goose.Up(db.Db, "../app/migrations")

}

func (t *ApiTest) TestCreatingNewUser() {
	data := url.Values{}
	data.Set("email", "1@test.com")
	data.Set("name", "Vardenis Pavardenis")
	data.Set("password", "slaptazodis")
	t.PostForm("/user/create", data)
	t.AssertStatus(201)

	var user models.User
	controllers.Dbm.SelectOne(&user, "SELECT * FROM users WHERE Email=$1", "1@test.com")
	t.AssertEqual("1@test.com", user.Email)
	t.AssertEqual(false, user.IsSuperuser)
	t.AssertEqual("Vardenis Pavardenis", user.Name)
}

func (t *ApiTest) TestUniqueEmailsOnly() {
	data := url.Values{}
	data.Set("email", "rugilena@gmail.com")
	data.Set("name", "Vardenis Pavardenis")
	data.Set("password", "slaptazodis")
	t.PostForm("/user/create", data)
	t.AssertStatus(400)
	fmt.Println(t.Response.Body)
}

func (t *ApiTest) After() {
	println("Tear down")
}
