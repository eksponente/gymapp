package tests

import (
	"fmt"

	db "github.com/revel/modules/db/app"
	"github.com/revel/revel/testing"

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
	db.Init()
}

func (t *AppTest) After() {
	println("Tear down")
}
