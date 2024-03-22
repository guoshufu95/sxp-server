package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	ini "sxp-server/common/initial"
	"sxp-server/common/logger"
	"sxp-server/common/model"
)

func Permission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var log *logger.ZapLog
		v, ok := ctx.Get("sxp-log")
		if !ok || v == nil {
			log = logger.GetLogger()
		} else {
			log = v.(*logger.ZapLog)
		}
		data, _ := ctx.Get("sxp-claims")
		claims := data.(*model.MyClaims)
		if claims.RoleKey == "admin" { //admin用户直接放行
			ctx.Next()
			return
		}
		e := ini.App.GetCasbin()
		if e == nil {
			log.Error("获取casbin失败")
			ctx.JSON(http.StatusOK, gin.H{
				"code": 403,
				"msg":  "casbin获取失败",
			})
			ctx.Abort()
			ctx.Abort()
			return
		}
		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		res, err := e.Enforce(claims.RoleKey, path, method)
		if err != nil {
			log.Errorf("casbin校验错误：%s", err.Error())
			return
		}
		if res {
			log.Info("casbin校验通过")
			ctx.Next()
		} else {
			log.Error("没有该接口访问权限")
			ctx.JSON(http.StatusOK, gin.H{
				"code": 403,
				"msg":  "没有该接口访问权限",
			})
			ctx.Abort()
			return
		}
	}
}
