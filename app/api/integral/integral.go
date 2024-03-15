package integral

import (
	"github.com/gin-gonic/gin"
	"sxp-server/app/api"
	"sxp-server/app/model"
	serv "sxp-server/app/service"
	"sxp-server/app/service/integral"
)

type IntegralApi struct {
	api.Api
}

var is integral.IntegralService

// InitIntegral
//
//	@Description: 初始化积分相关信息
//	@receiver a
//	@param c
func (a *IntegralApi) InitIntegral(c *gin.Context) {
	a.MakeApi(c)
	var req model.IntegralReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	serv.MakeService(&is.Service, c)
	req.RemainCount = req.Count
	req.RemainIntegral = req.Integral
	err = is.InitIntegral(req)
	if err != nil {
		a.Logger.Error("初始化积分失败!")
		a.ResponseError(err)
		return
	}
	a.Response("积分初始化成功!", nil)
}
