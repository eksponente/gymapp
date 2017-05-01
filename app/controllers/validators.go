package controllers

import (
	"regexp"

	valid "github.com/asaskevich/govalidator"
	"github.com/revel/revel"
)

type EmailValidator struct{}

func (v EmailValidator) IsSatisfied(email interface{}) bool {
	return valid.IsEmail(email.(string))
}

func (v EmailValidator) DefaultMessage() string {
	return "Invalid email."
}

type OverrideMesage struct {
	revel.Validator
	Message string
}

func (v OverrideMesage) DefaultMessage() string {
	return v.Message
}

type RegexpValidator struct {
	Message string
	Reg     string
}

func (v RegexpValidator) DefaultMessage() string {
	return v.Message
}

func (v RegexpValidator) IsSatisfied(inp interface{}) bool {
	return revel.Match{regexp.MustCompile(v.Reg)}.IsSatisfied(inp.(string))
}
