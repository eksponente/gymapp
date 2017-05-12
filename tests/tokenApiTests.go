package tests

import (
	"encoding/json"
	"errors"
	"gymapp/app/controllers"
	"gymapp/app/models"
	"net/url"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/pressly/goose"
	db "github.com/revel/modules/db/app"
	"github.com/revel/revel"
	"github.com/revel/revel/testing"
)

type TokenApiTest struct {
	testing.TestSuite
}

var user models.User

func (t *TokenApiTest) Before() {
	println("Set up")
	if name, _ := revel.Config.String("db.name"); name != "gymapptest" {
		panic(errors.New("Not connected to test database. RDS_DB_NAME must be gymapptest."))
	}
	println("MIGRATING DOWN")
	goose.DownTo(db.Db, "../app/migrations", 20170419150037)
	goose.Down(db.Db, "../app/migrations")
	println("MIGRATING UP")
	goose.Up(db.Db, "../app/migrations")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("slaptazodis"), bcrypt.DefaultCost)
	res, err := controllers.Dbm.Exec("insert into \"users\" (\"name\", \"email\", \"issuperuser\", \"password\") VALUES ($1, $2, $3, $4);", "Test User", "1@test.com", false, string(hashedPassword))

	if rows, _ := res.RowsAffected(); err != nil || rows == 0 {
		panic("Something wrong")
	}
	controllers.Dbm.SelectOne(&user, "SELECT * FROM users WHERE Email=$1", "1@test.com")
}

func (t *TokenApiTest) After() {
	println("Tear down")
}

//TOKEN API TESTS
func (t *TokenApiTest) TestCreatingNewToken() {
	data := url.Values{}
	data.Set("email", user.Email)
	data.Set("password", "slaptazodis")
	t.PostForm("/token/request", data)
	t.AssertStatus(201)

	var token models.Token
	controllers.Dbm.SelectOne(&token, "SELECT * FROM tokens WHERE user_id=$1", user.UserId)
	t.AssertContains(token.Token)
	t.AssertContains(token.ExpirationDate)
}

func (t *TokenApiTest) TestUnableToCreateTokenWithWrongPassword() {
	data := url.Values{}
	data.Set("email", user.Email)
	data.Set("password", "wrongPassword")
	t.PostForm("/token/request", data)
	t.AssertStatus(404)

	var token models.Token
	err := controllers.Dbm.SelectOne(&token, "SELECT * FROM tokens WHERE user_id=$1", user.UserId)
	if err == nil {
		panic(errors.New("A token has been created."))
	}
}

func (t *TokenApiTest) TestDestroyingToken() {
	data := url.Values{}
	data.Set("email", "1@test.com")
	data.Set("password", "slaptazodis")
	t.PostForm("/token/request", data)
	println(string(t.ResponseBody))
	t.AssertStatus(201)

	var resp map[string]interface{}
	json.Unmarshal(t.ResponseBody, &resp)
	token := resp["token"].(string)

	data = url.Values{}
	data.Set("token", token)
	t.PostForm("/token/destroy", data)
	t.AssertStatus(200)
	result, _ := controllers.Dbm.Exec("SELECT * FROM tokens WHERE user_id=$1", user.UserId)

	rows, _ := result.RowsAffected()
	t.AssertEqual(0, rows)
}

func (t *TokenApiTest) TestRenewingToken() {
	data := url.Values{}
	data.Set("email", "1@test.com")
	data.Set("password", "slaptazodis")
	t.PostForm("/token/request", data)
	t.AssertStatus(201)

	var resp map[string]interface{}
	json.Unmarshal(t.ResponseBody, &resp)
	token := resp["token"].(string)
	oldTokenExp, _ := time.ParseInLocation(time.RFC3339, resp["expiration"].(string), controllers.Location)

	time.Sleep(2 * time.Second)

	data = url.Values{}
	data.Set("token", token)

	var tok models.Token
	controllers.Dbm.SelectOne(&tok, "SELECT * FROM tokens WHERE user_id=$1", user.UserId)

	t.PostForm("/token/renew", data)
	t.AssertStatus(200)

	json.Unmarshal(t.ResponseBody, &resp)
	newTokenExp, _ := time.ParseInLocation(time.RFC3339, resp["expiration"].(string), controllers.Location)
	t.Assert(newTokenExp.After(oldTokenExp))
}

func (t *TokenApiTest) TestDatabase() {
	controllers.Dbm.Exec("INSERT INTO \"tokens\" (\"user_id\",\"token\",\"expirationdate\") VALUES ($1,$2,$3);", user.UserId, "tokenas", time.Now().Add(time.Hour*24*14).In(controllers.Location).Format(time.RFC3339))
	var tok models.Token
	controllers.Dbm.SelectOne(&tok, "SELECT * FROM tokens WHERE user_id=$1", user.UserId)
	println("ZIUREK CIA")
	println(tok.Token)
	t.AssertEqual("tokenas", tok.Token)
}
