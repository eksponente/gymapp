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

type Token struct {
	GorpController
}

type User struct {
	GorpController
}

//InitDB initializes the database for the application usages
func InitDB() {
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

func RetrieveUser(email string, c GorpController) (user models.User, err error) {
	err = c.Txn.SelectOne(&user, "SELECT * FROM users WHERE Email=$1", email)
	return
}

func RetrieveToken(t string, c GorpController) (token models.Token, err error) {
	err = c.Txn.SelectOne(&token, "SELECT * FROM tokens WHERE Token=$1", t)
	return
}

func UpdateTokenExpDate(t string, exp string, c GorpController) (err error) {
	stmt, err := c.Txn.Prepare("UPDATE \"tokens\" SET \"expirationdate\" = $1 WHERE \"token\" = '$2';")
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

func CreateUser(name, email, password string, c GorpController) (err error, rows int64) {
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
