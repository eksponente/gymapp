package controllers

import (
	"fmt"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/revel/revel"

	"golang.org/x/crypto/bcrypt"
)

//RequestToken request a token for login
func (c App) RequestToken() revel.Result {
	username := c.Params.Form.Get("username")
	password := c.Params.Form.Get("password")

	c.Validation.Match(username, regexp.MustCompile("^\\w*$")).Message("Username can only consist letter characters.")

	//Check if username and password valid
	user, err1 := RetrieveUser(username, c)
	err2 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err1 != nil || err2 != nil {
		m := make(map[string]string)
		m["error"] = "Invalid username or password."
		return c.RenderJSON(m)
	}

	//Create a new token and save it in database
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24 * 14).Format(time.RFC3339)
	token.Claims = claims
	secret, _ := revel.Config.String("secret")
	signedToken, _ := token.SignedString([]byte(secret))
	fmt.Println(signedToken)
	CreateToken(string(signedToken), user, claims["exp"].(string), c)

	//return the token
	m := make(map[string]string)
	m["token"] = signedToken
	m["expiration"] = claims["exp"].(string)
	return c.RenderJSON(m)
}

func (c App) RenewToken() revel.Result {
	t := c.Params.Form.Get("token")
	_, err := RetrieveToken(t, c)
	if err != nil {
		m := make(map[string]string)
		m["error"] = "Token not found."
		return c.RenderJSON(m)
	}
	exp := time.Now().Add(time.Hour * 24 * 14).Format(time.RFC3339)
	UpdateTokenExpDate(t, exp, c)

	//return the token
	m := make(map[string]string)
	m["token"] = t
	m["expiration"] = exp
	return c.RenderJSON(m)
}
