package controllers

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

//Index test creating user
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

	//TODO: validate those inputs

	if !valid.IsEmail(email) {
		m := make(map[string]string)
		c.Response.Status = 400
		m["error"] = "Invalid email."
		return c.RenderJSON(m)
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	err, rows := CreateUser(name, email, string(hashedPassword), c.GorpController)
	if rows == 0 { //no rows have been created
		m := make(map[string]string)
		m["message"] = "User with that email already exists."
		c.Response.Status = 400
		return c.RenderJSON(m)
	}
	if err != nil {
		m := make(map[string]string)
		m["error"] = "Database error."
		c.Response.Status = 400
		return c.RenderJSON(m)
	}

	m := map[string]interface{}{
		"error": nil,
	}
	c.Response.Status = 201
	return c.RenderJSON(m)
}
