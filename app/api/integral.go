package api

import (
	"github.com/gin-gonic/gin"
	serv "sxp-server/app/service"
	"sxp-server/app/service/dto"
)

type IntegralApi struct {
	Api
}

var is serv.IntegralService

// InitIntegral
//
//	@Description: 初始化积分相关信息
//	@receiver a
//	@param c
func (a *IntegralApi) InitIntegral(c *gin.Context) {
	a.BuildApi(c).BuildService(&is.Service)
	var req dto.IntegralReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
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

// DoIntegral
//
//	@Description: 抢积分/红包
//	@receiver a
//	@param c
func (a *IntegralApi) DoIntegral(c *gin.Context) {
	a.BuildApi(c).BuildService(&is.Service)
	var req dto.DoIntegral
	//todo 用户校验相关的内容省略，不是重点
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	err = is.Do()
}
