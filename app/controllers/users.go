package controllers

import "github.com/revel/revel"

//Index test creating user
func (c App) createUser() revel.Result {
	return c.Render()
}
