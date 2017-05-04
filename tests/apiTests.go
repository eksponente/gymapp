package tests

import (
	"encoding/json"
	"errors"
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
	data.Set("name", "Vardenis Pa-Vardenis")
	data.Set("password", "slaptazodis")
	t.PostForm("/user/create", data)
	t.AssertStatus(201)

	var user models.User
	controllers.Dbm.SelectOne(&user, "SELECT * FROM users WHERE Email=$1", "1@test.com")
	t.AssertEqual("1@test.com", user.Email)
	t.AssertEqual(false, user.IsSuperuser)
	t.AssertEqual("Vardenis Pa-Vardenis", user.Name)
}

func (t *ApiTest) TestUniqueEmailsOnly() {
	data := url.Values{}
	data.Set("email", "rugilena@gmail.com")
	data.Set("name", "Vardenis Pavardenis")
	data.Set("password", "slaptazodis")
	t.PostForm("/user/create", data)
	t.AssertStatus(400)
	var res map[string]interface{}
	json.Unmarshal(t.ResponseBody, &res)
	t.AssertEqual("User with that email already exists.", res["error"].(string))

}

func (t *ApiTest) TestValidEmailsOnly() {
	data := url.Values{}
	data.Set("email", "rugilena@gmail")
	data.Set("name", "Vardenis Pavardenis")
	data.Set("password", "slaptazodis")
	t.PostForm("/user/create", data)
	t.AssertStatus(400)

	var user models.User
	err := controllers.Dbm.SelectOne(&user, "SELECT * FROM users WHERE Email=$1", "rugilena@gmail")
	if err == nil {
		panic(errors.New("New user has been created."))
	}

	t.AssertContains("Invalid email.")

	// t.AssertEqual(res["errors"].([]interface{})[0].(map[string]interface{})["Message"].(string), "Invalid email.")
}

func (t *ApiTest) TestValidNamesOnly() {
	data := url.Values{}
	data.Set("email", "2@test.com")
	data.Set("name", "()()()")
	data.Set("password", "slaptazodis")
	t.PostForm("/user/create", data)
	t.AssertStatus(400)

	var user models.User
	err := controllers.Dbm.SelectOne(&user, "SELECT * FROM users WHERE Email=$1", "2@test.com")
	if err == nil {
		panic(errors.New("New user has been created."))
	}

	// t.AssertEqual(res["errors"].([]interface{})[0].(map[string]interface{})["Message"].(string), "Valid name is required.")
	t.AssertContains("Valid name is required.")
}

func (t *ApiTest) TestValidPasswordsOnly() {
	data := url.Values{}
	data.Set("email", "3@test.com")
	data.Set("name", "Vardas")
	data.Set("password", "sla")
	t.PostForm("/user/create", data)
	t.AssertStatus(400)

	var user models.User
	err := controllers.Dbm.SelectOne(&user, "SELECT * FROM users WHERE Email=$1", "3@test.com")
	if err == nil {
		panic(errors.New("New user has been created."))
	}

	t.AssertContains("Minimum size is 8")
}

func (t *ApiTest) After() {
	println("Tear down")
}
