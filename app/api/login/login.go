package login

import (
	"github.com/gin-gonic/gin"
	"sxp-server/app/api"
	"sxp-server/app/model"
	serv "sxp-server/app/service"
	"sxp-server/app/service/login"
	"sxp-server/app/service/task"
)

type LoginApi struct {
	api.Api
}

var ts = task.TaskService{}

//func init() {
//	serv.MakeService(&ts.Service)
//}

// Login
//
//	@Description: 登录
//	@receiver l
//	@param c
func (l *LoginApi) Login(c *gin.Context) {
	l.MakeApi(c)
	var req model.LoginReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		l.Logger.Error(err.Error())
		l.ResponseError(err)
		return
	}
	var s login.LoginService
	serv.MakeService(&s.Service, c)
	err, token := s.Login(req)
	if err != nil {
		l.Logger.Error(err.Error())
		l.ResponseError(err)
	}
	l.Response("成功生成token", token)
}
