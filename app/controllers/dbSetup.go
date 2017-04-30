package controllers

import (
	"database/sql"
	"fmt"
	"gymapp/app/models"
	"time"

	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq" //postgres driver
	"github.com/pressly/goose"

	"github.com/revel/modules/db/app"
	r "github.com/revel/revel"
)

var (
	//Dbm database mapping
	Dbm *gorp.DbMap
	//Location used for timezone parsing
	Location *time.Location
)

//App is a GorpController wrapper for App endpoints
type App struct {
	GorpController
}

//Token is a GorpController wrapper for token api endpoints
type Token struct {
	GorpController
}

//User is a GorpController wrapper for User endpoints
type User struct {
	GorpController
}

//InitDB initializes the database for the application usages
func InitDB() {
	fmt.Print(r.Config.String("db.spec"))
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.PostgresDialect{}}
	var err error
	Location, err = time.LoadLocation(r.Config.StringDefault("location", "Europe/London"))
	if err != nil {
		panic(err)
	}

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

//In the following section we will define functions which will take care of database calls needed

//RetrieveUser will retrieve a user from the database
func RetrieveUser(email string, c GorpController) (user models.User, err error) {
	err = c.Txn.SelectOne(&user, "SELECT * FROM users WHERE Email=$1", email)
	return
}

//RetrieveToken will retrieve a token from the database
func RetrieveToken(t string, c GorpController) (token models.Token, err error) {
	err = c.Txn.SelectOne(&token, "SELECT * FROM tokens WHERE Token=$1", t)
	return
}

//UpdateTokenExpDate will update a token expiration date
func UpdateTokenExpDate(t string, exp string, c GorpController) (err error) {
	stmt, err := c.Txn.Prepare("UPDATE \"tokens\" SET \"expirationdate\" = $1 WHERE \"token\" = $2;")
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

//DeleteToken will delete a token
func DeleteToken(t string, c GorpController) (err error) {
	stmt, err := c.Txn.Prepare("DELETE FROM \"tokens\" WHERE \"token\" = $1;")
	if err != nil {
		return
	}
	defer stmt.Close()
	if _, err = stmt.Exec(t); err != nil {
		c.Txn.Rollback()
		return
	}
	return
}

//CreateToken will create a token
func CreateToken(t string, user models.User, exp string, c GorpController) (err error) {
	stmt, err := c.Txn.Prepare("insert into \"tokens\" (\"user_id\",\"token\",\"expirationdate\") values ($1,$2,$3);")
	if err != nil {
		return
	}
	defer stmt.Close()
	if _, err = stmt.Exec(user.UserId, t, exp); err != nil {
		c.Txn.Rollback()
		return
	}
	return
}

//CreateUser will create a user
func CreateUser(name, email, password string, c GorpController) (rows int64, err error) {
	stmt, err := c.Txn.Prepare("insert into \"users\" (\"name\", \"email\", \"issuperuser\", \"password\") VALUES ($1, $2, $3, $4);")
	if err != nil {
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(name, email, false, password)
	if err != nil {
		c.Txn.Rollback()
		return
	}
	rows, _ = res.RowsAffected()

	return
}
