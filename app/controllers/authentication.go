package controllers

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/revel/revel"

	valid "github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
)

//Request api endpoint to request a new token.
func (c Token) Request() revel.Result {
	email := c.Params.Form.Get("email")
	password := c.Params.Form.Get("password")

	//check if email is in a valid format
	if !valid.IsEmail(email) {
		m := make(map[string]string)
		c.Response.Status = 400
		m["error"] = "Invalid email."
		return c.RenderJSON(m)
	}

	//Check if username and password valid
	user, err1 := RetrieveUser(email, c.GorpController)
	err2 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err1 != nil || err2 != nil {
		m := make(map[string]string)
		c.Response.Status = 404
		m["error"] = "Invalid email or password."
		return c.RenderJSON(m)
	}

	//Create a new token and save it in database
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["email"] = email
	claims["exp"] = time.Now().In(Location).Add(time.Hour * 24 * 14).Format(time.RFC3339)
	token.Claims = claims
	secret, _ := revel.Config.String("secret")
	signedToken, _ := token.SignedString([]byte(secret))
	fmt.Println(signedToken)
	CreateToken(string(signedToken), user, claims["exp"].(string), c.GorpController)

	//return the token
	m := map[string]interface{}{
		"token":      signedToken,
		"expiration": claims["exp"],
		"error":      nil,
	}
	c.Response.Status = 200
	return c.RenderJSON(m)
}

//Renew an api endpoint to renew a certain token to be valid up until 2 weeks from now.
func (c Token) Renew() revel.Result {
	t := c.Params.Form.Get("token")
	token, errTok := RetrieveToken(t, c.GorpController)

	tokenExp, errTime := time.ParseInLocation(time.RFC3339, token.ExpirationDate, Location)

	if errTime != nil {
		panic(errTime)
		// c.Response.Status = 500
		// m := map[string]interface{}{
		// 	"error": "Internal server error.",
		// }
		// return c.RenderJSON(m)
	}

	if errTok != nil || time.Now().In(Location).After(tokenExp) {
		m := make(map[string]string)
		c.Response.Status = 404
		m["error"] = "Token not found or already expired."
		return c.RenderJSON(m)
	}

	exp := time.Now().Add(time.Hour * 24 * 14).Format(time.RFC3339)
	err := UpdateTokenExpDate(t, exp, c.GorpController)
	if err != nil {
		panic(err)
		// m := make(map[string]string)
		// m["error"] = "Internal server error."
		// c.Response.Status = 500
		// return c.RenderJSON(m)
	}

	//return the token
	m := map[string]interface{}{
		"token":      t,
		"expiration": exp,
		"error":      nil,
	}
	c.Response.Status = 200
	return c.RenderJSON(m)
}

//Destroy an api endpoint to destroy a token.
func (c Token) Destroy() revel.Result {
	t := c.Params.Form.Get("token")
	_, err := RetrieveToken(t, c.GorpController)
	if err != nil {
		m := make(map[string]string)
		m["error"] = "Token not found."
		c.Response.Status = 404
		return c.RenderJSON(m)
	}
	err = DeleteToken(t, c.GorpController)
	if err != nil {
		panic(err)
		// m := make(map[string]string)
		// m["error"] = "Internal server error."
		// c.Response.Status = 500
		// return c.RenderJSON(m)
	}

	//return the token
	m := map[string]interface{}{
		"error": nil,
	}
	c.Response.Status = 200
	return c.RenderJSON(m)
}
