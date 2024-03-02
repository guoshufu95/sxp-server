package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	ini "sxp-server/common/initial"
	"sxp-server/common/logger"
)

type Api struct {
	Logger *logger.ZapLog
	Ctx    *gin.Context
}

func (a *Api) MakeApi(c *gin.Context) {
	app := ini.MakeApp()
	a.Logger = app.Logger
	a.Ctx = c
}

func (a *Api) ResponseError(err error) {
	res := gin.H{
		"code":    http.StatusInternalServerError,
		"message": err.Error(),
	}
	a.Ctx.Set("status", http.StatusInternalServerError)
	a.Ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)

}

func (a *Api) Response(msg string, data ...interface{}) {
	res := gin.H{
		"code":    http.StatusOK,
		"message": msg,
		"data":    data,
	}
	a.Ctx.AbortWithStatusJSON(http.StatusOK, res)
	a.Ctx.Set("status", http.StatusOK)
}
