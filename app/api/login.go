package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	serv "sxp-server/app/service"
	"sxp-server/app/service/dto"
)

type LoginApi struct {
	Api
}

var ls = serv.LoginService{}

// Login
//
//	@Description: 登录
//	@receiver l
//	@param c
func (a LoginApi) Login(c *gin.Context) {
	a.BuildApi(c).BuildService(&ls.Service)
	var req dto.LoginReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.ResponseError(http.StatusBadRequest, err)
		return
	}
	err, token := ls.Login(req)
	if err != nil {
		a.ResponseError(http.StatusInternalServerError, err)
		return
	}
	a.Response("成功生成token", token)
}
