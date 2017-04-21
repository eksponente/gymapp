package controllers

import (
	"database/sql"
	"fmt"
	"gymapp/app/models"

	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq" //postgres driver
	"github.com/pressly/goose"

	"github.com/revel/modules/db/app"
	r "github.com/revel/revel"
)

var (
	//Dbm database mapping
	Dbm *gorp.DbMap
)

type App struct {
	GorpController
}

//InitDB initializes the database for the application usages
func InitDB() {
	// TODO: set the db.spec as environment variable as per http://revel.github.io/manual/appconf.html
	fmt.Print(r.Config.String("db.spec"))
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.PostgresDialect{}}

	//run the migrations
	goose.Run("up", db.Db, "../migrations")

	//set up gorp with the databse
	Dbm.AddTableWithName(models.User{}, "users").SetKeys(true, "user_id")
	Dbm.AddTableWithName(models.Token{}, "tokens")

	Dbm.TraceOn("[gorp]", r.INFO)
}

//GorpController custom controllers
type GorpController struct {
	*r.Controller
	Txn *gorp.Transaction
}

//Begin interceptor
func (c *GorpController) Begin() r.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

//Commit interceptor
func (c *GorpController) Commit() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

//Rollback interceptor
func (c *GorpController) Rollback() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func RetrieveUser(username string, c App) (user models.User, err error) {
	err = c.Txn.SelectOne(&user, "SELECT * FROM users WHERE Username=$1", username)
	return
}

func RetrieveToken(t string, c App) (token models.Token, err error) {
	err = c.Txn.SelectOne(&token, "SELECT * FROM tokens WHERE Token=$1", t)
	return
}

func UpdateTokenExpDate(t string, exp string, c App) (err error) {
	stmt, err := c.Txn.Prepare("UPDATE tokens SET expirationdate = '$1' WHERE token = '$2';")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(exp, t); err != nil {
		c.Txn.Rollback()
		return err
	}
	return err
}

func CreateToken(t string, user models.User, exp string, c App) (err error) {
	stmt, err := c.Txn.Prepare("insert into \"tokens\" (\"user_id\",\"token\",\"expirationdate\") values ($1,$2,$3);")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(user.UserId, t, exp); err != nil {
		c.Txn.Rollback()
		return err
	}
	return err
}
