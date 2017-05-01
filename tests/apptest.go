package tests

import (
	"fmt"

	"github.com/revel/revel/testing"

	"gymapp/app/models"

	r "github.com/revel/revel"
)

type AppTest struct {
	testing.TestSuite
}

func (t *AppTest) Before() {
	println("Set up")
}

func (t *AppTest) TestThatIndexPageWorks() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *AppTest) TestDatabaseConnection() {
	fmt.Print(r.Config.String("db.spec"))
	var user models.User
	Dbm.SelectOne(&user, "SELECT * FROM users WHERE Email=$1", "rugilena@gmail.com")
	t.AssertEqual(user.IsSuperuser, true)

}

func (t *AppTest) After() {
	println("Tear down")
}
