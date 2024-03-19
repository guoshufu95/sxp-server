package model

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type MyClaims struct {
	Username string `json:"username"`
	RoleId   int    `json:"roleId"`
	RoleKey  string `json:"roleKey"`
	jwt.StandardClaims
}
