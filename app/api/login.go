package api

import (
	"github.com/gin-gonic/gin"
	serv "sxp-server/app/service"
	"sxp-server/app/service/dto"
)

type LoginApi struct {
	Api
}

var ls = serv.LoginService{}

//func init() {
//	serv.MakeService(&ts.Service)
//}

// Login
//
//	@Description: 登录
//	@receiver l
//	@param c
func (a *LoginApi) Login(c *gin.Context) {
	a.BuildApi(c).BuildService(&ls.Service)
	var req dto.LoginReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	//serv.MakeService(&ls.Service, c)
	err, token := ls.Login(req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	a.Response("成功生成token", token)
}
