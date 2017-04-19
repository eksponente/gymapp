package controllers

import (
	"gymapp/app/models"

	"github.com/revel/revel"
)

//App controller
type App struct {
	GorpController
}

//Index test creating user
func (c App) Index() revel.Result {
	err := c.Txn.Insert(&models.User{Name: "Rugile", Password: "maironis"})
	if err != nil {
		panic(err)
	}
	return c.Render()
}
