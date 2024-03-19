package jwtToken

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"sxp-server/common/model"
	"sxp-server/config"
	"time"
)

//var SECRETKEY = []byte("sxp-server") //私钥

// GenToken 生成JWT
func GenToken(username, roleKey string, roleId int) (string, error) {
	expTime := time.Now().Add(time.Duration(config.Conf.Jwt.Timeout) * time.Second).Unix()
	c := model.MyClaims{
		Username: username,
		RoleId:   roleId,
		RoleKey:  roleKey,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime,      // 过期时间
			Issuer:    "sxp-server", // 签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(config.Conf.Jwt.Secret))
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*model.MyClaims, *jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(config.Conf.Jwt.Secret), nil
	})
	if err != nil {
		return nil, nil, err
	}
	if claims, ok := token.Claims.(*model.MyClaims); ok && token.Valid { // 校验token
		return claims, token, nil
	}
	return nil, nil, errors.New("invalid token")
}
