package controllers

import (
	"database/sql"
	"gymapp/app/models"
	"strconv"
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
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.PostgresDialect{}}
	var err error
	Location, err = time.LoadLocation(r.Config.StringDefault("location", "Europe/London"))
	if err != nil {
		panic(err)
	}

	//run the migrations
	goose.Up(db.Db, "../migrations")

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
func RetrieveUser(email string, txn *gorp.Transaction) (user models.User, err error) {
	err = txn.SelectOne(&user, "SELECT * FROM \"users\" WHERE \"email\"=$1;", email)
	return
}

//RetrieveToken will retrieve a token from the database
func RetrieveToken(t string, txn *gorp.Transaction) (token models.Token, err error) {
	err = txn.SelectOne(&token, "SELECT * FROM \"tokens\" WHERE \"token\"=$1;", t)
	return
}

//UpdateTokenExpDate will update a token expiration date
func UpdateTokenExpDate(t string, exp string, txn *gorp.Transaction) (err error) {
	stmt, err := txn.Prepare("UPDATE \"tokens\" SET \"expirationdate\" = $1 WHERE \"token\" = $2;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(exp, t); err != nil {
		txn.Rollback()
		return err
	}
	return err
}

//DeleteToken will delete a token
func DeleteToken(t string, txn *gorp.Transaction) (err error) {
	stmt, err := txn.Prepare("DELETE FROM \"tokens\" WHERE \"token\" = $1;")
	if err != nil {
		return
	}
	defer stmt.Close()
	if _, err = stmt.Exec(t); err != nil {
		txn.Rollback()
		return
	}
	return
}

//CreateToken will create a token
func CreateToken(t string, user models.User, exp string, txn *gorp.Transaction) (err error) {
	stmt, err := txn.Prepare("insert into \"tokens\" (\"user_id\",\"token\",\"expirationdate\") values ($1,$2,$3);")
	if err != nil {
		return
	}
	defer stmt.Close()
	if _, err = stmt.Exec(user.UserId, t, exp); err != nil {
		txn.Rollback()
		return
	}
	return
}

//CreateUser will create a user
func CreateUser(name, email, password string, txn *gorp.Transaction) (rows int64, err error) {
	stmt, err := txn.Prepare("insert into \"users\" (\"name\", \"email\", \"issuperuser\", \"password\") VALUES ($1, $2, $3, $4);")
	if err != nil {
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(name, email, false, password)
	if err != nil {
		txn.Rollback()
		return
	}
	rows, _ = res.RowsAffected()

	return
}

//CreateLocation will create a user
func CreateLocation(user_id int, latitude float64, longitude float64, address string, location_name string, txn *gorp.Transaction) (rows int64, err error) {
	stmt, err := txn.Prepare("INSERT INTO \"locations\" (\"coordinates\", \"address\", \"location_name\", \"user_id\", \"lat\", \"lon\") values (ST_GeomFromText($1, 4326), $2, $3, $4, $5, $6);")
	if err != nil {
		return
	}
	defer stmt.Close()
	geometryPoint := "POINT(" + strconv.FormatFloat(latitude, 'f', -1, 64) + " " + strconv.FormatFloat(longitude, 'f', -1, 64) + ")"
	res, err := stmt.Exec(geometryPoint, address, location_name, user_id, latitude, longitude)
	if err != nil {
		txn.Rollback()
		return
	}
	rows, _ = res.RowsAffected()

	return
}

// test_geom.coord <-> ST_GeographyFromText('SRID=4326;POINT(-100 50)') limit 1;
