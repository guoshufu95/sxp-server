package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sxp-server/common/jwtToken"
	"sxp-server/common/logger"
)

// JWTAuthMiddleware
//
//	@Description: 基于JWT的认证中间件
//	@return gin.HandlerFunc
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.GetLogger()
		header := c.Request.Header.Get("Authorization")
		if header == "" {
			err := errors.New("请传入合法token")
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  http.StatusUnauthorized,
				"error": err.Error()})
			return
		}
		parts := strings.SplitN(header, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			err := errors.New("请求头中auth格式有误")
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  http.StatusUnauthorized,
				"error": err.Error()})
			return
		}
		mc, err := jwtToken.ParseToken(parts[1])
		if err != nil {
			err = errors.New("token解析失败！")
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  http.StatusUnauthorized,
				"error": err.Error()})
			return
		}
		c.Set("username", mc.Username)
		c.Next()
	}
}
