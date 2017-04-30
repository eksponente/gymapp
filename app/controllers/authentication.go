// @SubApi User API [/token]
package controllers

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/revel/revel"

	valid "github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
)

// @Title Request
// @Description Request a token
// @Param   email     form    string     true        "User email"
// @Param   password     form    string     true        "User password"
// @Success 200 {object} string
// @Failure 404 {object} APIError "Invalid email or password."
// @Router /token/request/ [post]
func (c Token) Request() revel.Result {
	email := c.Params.Form.Get("email")
	password := c.Params.Form.Get("password")

	if !valid.IsEmail(email) {
		// Store the validation errors in the flash context and redirect.
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
	claims["exp"] = time.Now().Add(time.Hour * 24 * 14).Format(time.RFC3339)
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

//RenewToken can be used to renew a certain token to be valid up until 2 weeks from now.
func (c Token) Renew() revel.Result {
	t := c.Params.Form.Get("token")
	_, err := RetrieveToken(t, c.GorpController)
	if err != nil {
		m := make(map[string]string)
		c.Response.Status = 404
		m["error"] = "Token not found."
		return c.RenderJSON(m)
	}
	exp := time.Now().Add(time.Hour * 24 * 14).Format(time.RFC3339)
	err = UpdateTokenExpDate(t, exp, c.GorpController)
	if err != nil {
		m := make(map[string]string)
		m["error"] = "Database error."
		c.Response.Status = 400
		return c.RenderJSON(m)
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

//DestroyToken can be used to renew a certain token to be valid up until 2 weeks from now.
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
		m := make(map[string]string)
		m["error"] = "Database error."
		c.Response.Status = 400
		return c.RenderJSON(m)
	}

	//return the token
	m := map[string]interface{}{
		"error": nil,
	}
	c.Response.Status = 200
	return c.RenderJSON(m)
}
