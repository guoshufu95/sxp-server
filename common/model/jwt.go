package model

import (
	jwt "github.com/dgrijalva/jwt-go"
)

// MyClaims
// @Description: 生成token用
type MyClaims struct {
	Username string   `json:"username"`
	RoleIds  []int    `json:"roleIds"`
	RoleKeys []string `json:"roleKeys"`
	jwt.StandardClaims
}
