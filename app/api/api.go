package api

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"sxp-server/app/service"
	_const "sxp-server/common/const"
	"sxp-server/common/logger"
)

type Api struct {
	Logger *logger.ZapLog
	Ctx    *gin.Context
}

// BuildApi
//
//	@Description: 初始api的一些字段
//	@receiver a
//	@param c
func (a *Api) BuildApi(c *gin.Context) *Api {
	a.Logger = c.MustGet(_const.SxpLogKey).(*logger.ZapLog)
	a.Ctx = c
	return a
}

// BuildService
//
//	@Description: 初始化service的一些字段
//	@receiver a
//	@param s
func (a *Api) BuildService(s *service.Service) {
	s.Db = a.Ctx.MustGet(_const.SxpGormDBkEY).(*gorm.DB)
	s.Cache = a.Ctx.MustGet(_const.SxpRedisDbKey).(*redis.Client)
	s.Logger = a.Logger
	s.Db.WithContext(a.Ctx)
}

// ResponseError
//
//	@Description: 错误返回
//	@receiver a
//	@param err
func (a *Api) ResponseError(code int, err error) {
	a.Logger.Error(err.Error())
	res := gin.H{
		"code":    code,
		"message": err.Error(),
	}
	a.Ctx.Set("status", code)
	a.Ctx.JSON(http.StatusOK, res)

}

// Response
//
//	@Description: 正常返回
//	@receiver a
//	@param msg
//	@param data
func (a *Api) Response(msg string, data ...interface{}) {
	a.Logger.Info("返回成功!")
	res := gin.H{
		"code":    http.StatusOK,
		"message": msg,
		"data":    data,
	}
	a.Ctx.Set("response", res)
	a.Ctx.Set("status", http.StatusOK)
	a.Ctx.JSON(http.StatusOK, res)
}
