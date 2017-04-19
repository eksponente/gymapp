package controllers

import (
	"database/sql"
	"fmt"
	"gymapp/app/models"

	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"

	"github.com/revel/modules/db/app"
	r "github.com/revel/revel"
)

var (
	Dbm *gorp.DbMap
)

func InitDB() {
	fmt.Print(r.Config.String("db.spec"))
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.PostgresDialect{}}

	goose.Run("up", db.Db, "../migrations")
	Dbm.AddTableWithName(models.User{}, "users").SetKeys(true, "user_id")

	Dbm.TraceOn("[gorp]", r.INFO)
}

type GorpController struct {
	*r.Controller
	Txn *gorp.Transaction
}

func (c *GorpController) Begin() r.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

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
