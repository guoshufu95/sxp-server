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
func GenToken(username string) (string, error) {
	expTime := time.Now().Add(time.Duration(config.Conf.Jwt.Timeout) * time.Second).Unix()
	c := model.MyClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime,      // 过期时间
			Issuer:    "sxp-server", // 签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(config.Conf.Jwt.Secret))
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*model.MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(config.Conf.Jwt.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*model.MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
