package tests

import (
	"github.com/revel/revel/testing"

	"gymapp/app/controllers"
	"gymapp/app/models"
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
	var user models.User
	controllers.Dbm.SelectOne(&user, "SELECT * FROM users WHERE Email=$1", "rugilena@gmail.com")
	t.AssertEqual(user.IsSuperuser, true)

}

func (t *AppTest) After() {
	println("Tear down")
}
