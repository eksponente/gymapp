package controllers

import (
	"regexp"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

//Create an api endpoint to create a user
func (c User) Create() revel.Result {
	email := c.Params.Form.Get("email")
	password := c.Params.Form.Get("password")
	name := c.Params.Form.Get("name")

	//validate all inputs
	c.Validation.Check(email, revel.Required{}, EmailValidator{})
	c.Validation.Check(name, revel.Required{}, OverrideMesage{revel.Match{regexp.MustCompile("^[\\p{L}\\s'.-]+$")}, "Valid name is required."})
	c.Validation.Check(password, revel.MinSize{8}, revel.Required{}, RegexpValidator{"Password must only contain latin characters and numbers", "^[a-zA-Z0-9]+$"})

	if c.Validation.HasErrors() {
		m := map[string]interface{}{
			"errors": c.Validation.Errors,
		}
		c.Response.Status = 400
		return c.RenderJSON(m)
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	rows, err := CreateUser(name, email, string(hashedPassword), c.GorpController.Txn)
	if rows == 0 { //no rows have been created
		m := make(map[string]string)
		m["error"] = "User with that email already exists."
		c.Response.Status = 400
		return c.RenderJSON(m)
	}
	if err != nil {
		panic(err)
	}

	m := map[string]interface{}{
		"error": nil,
	}
	c.Response.Status = 201
	return c.RenderJSON(m)
}
