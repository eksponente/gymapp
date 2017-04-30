package controllers

import (
	"regexp"

	valid "github.com/asaskevich/govalidator"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

//Create an api endpoint to create a user
func (c User) Create() revel.Result {
	email := c.Params.Form.Get("email")
	password := c.Params.Form.Get("password")
	name := c.Params.Form.Get("name")

	if email == "" || name == "" || password == "" {
		m := make(map[string]string)
		c.Response.Status = 400
		m["error"] = "Required fields: email, name, password."
		return c.RenderJSON(m)
	}

	//validate all inputs

	//validate email format
	if !valid.IsEmail(email) {
		m := make(map[string]string)
		c.Response.Status = 400
		m["error"] = "Invalid email."
		return c.RenderJSON(m)
	}

	//validate password (at least 8 chars long, letters and numbers only)
	c.Validation.Match(password, regexp.MustCompile("^[a-zA-Z0-9]+$")).Message("Password must only contain uppercase and lowercase latin characters or numbers.")
	c.Validation.MinSize(password, 8).Message("Password must be at least 8 characters long.")

	//validate name (unicode and some punctuation characters)
	c.Validation.Match(name, regexp.MustCompile("^[\\p{L}\\s'.-]+$")).Message("Name can only contain unicode letters and symbols '.-")

	if c.Validation.HasErrors() {
		m := map[string]interface{}{
			"errors": c.Validation.Errors,
		}
		c.Response.Status = 400
		return c.RenderJSON(m)
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	rows, err := CreateUser(name, email, string(hashedPassword), c.GorpController)
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
