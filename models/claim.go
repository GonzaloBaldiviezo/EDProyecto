package models

import (
	jwt "github.com/dgrijalva/jwt-go"
)

//Claim token de usuario
type Claim struct {
	User `json:"user"`
	jwt.StandardClaims
}
