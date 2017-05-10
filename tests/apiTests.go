package tests

import (
	"github.com/revel/revel/testing"
)

type ApiTest struct {
	testing.TestSuite
}

func (t *ApiTest) Before() {
	println("Set up")
}

func (t *ApiTest) TestThatIndexPageWorks() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *ApiTest) TestCreatingNewUser() {
    data := url.Values{
        "email": "example@user.com"
        "name": "Vardenis Pavardenis"
        "password": "kcromuva"
    }
    t.PostForm("/user/create", data)

}

func (t *ApiTest) After() {
	println("Tear down")
}
