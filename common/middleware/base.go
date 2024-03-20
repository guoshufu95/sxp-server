package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	_const "sxp-server/common/const"
	"time"
)

// WithGormDb
//
//	@Description: db中间件
//	@param db
//	@return gin.HandlerFunc
func WithGormDb(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(_const.SxpGormDBkEY, db.WithContext(ctx.Copy()))
		ctx.Next()
	}
}

// WithRedisDb
//
//	@Description: redis中间件
//	@param rdb
//	@return gin.HandlerFunc
func WithRedisDb(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(_const.SxpRedisDbKey, rdb)
		ctx.Next()
	}
}

// CORS 跨域请求
func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "POST", "GET", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 24 * time.Hour,
	})
}
