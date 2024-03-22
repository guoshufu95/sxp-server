package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"sxp-server/app/dao"
	_const "sxp-server/common/const"
	"sxp-server/common/logger"
	cm "sxp-server/common/model"
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

// Auth
//
//	@Description: 权限检查
//	@param db
//	@param c
//	@return err
//	@return flag
func (s *Service) Auth(c *gin.Context) (err error) {
	v, ok := c.Get(_const.SxpClaimsKey)
	if !ok {
		err = errors.New("无法获取claims")
		return
	}
	claims := v.(*cm.MyClaims)
	err, user := dao.GetAuth(s.Db, claims.RoleId)
	if err != nil {
		err = errors.New("获取当前登录用户信息失败")
		return
	}
	if user.IsSuper == 0 {
		err = errors.New("权限不足")
		return
	}
	return
}
