package service

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"sxp-server/common/logger"
)

type Service struct {
	Db     *gorm.DB
	Cache  *redis.Client
	Logger *logger.ZapLog
}

// MakeService
//
//	@Description: 初始化的一些赋值
//	@param s
//	@param c
func MakeService(s *Service, c *gin.Context) {
	s.Db = c.MustGet("sxp_gorm_db").(*gorm.DB)
	s.Cache = c.MustGet("sxp_redis_db").(*redis.Client)
	s.Logger = c.MustGet("sxp_zap_log").(*logger.ZapLog)
}
