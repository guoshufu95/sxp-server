package login

import (
	"sxp-server/app/dao/login"
	"sxp-server/app/model"
	"sxp-server/app/service"
	"sxp-server/common/jwtToken"
	"sxp-server/common/utils"
)

type LoginService struct {
	service.Service
}

// Login
//
//	@Description: 登录校验相关
//	@receiver s
//	@param req
//	@return err
func (s *LoginService) Login(req model.LoginReq) (err error, token string) {
	err, user := login.GetUser(s.Db, req.Username)
	if err != nil {
		s.Logger.Error("根据用户名查询用户失败!")
		return
	}
	_, err = utils.CompareHashAndPassword(user.PassWord, req.Password)
	if err != nil {
		s.Logger.Error("用户密码不匹配")
		return
	}
	token, err = jwtToken.GenToken(req.Username)
	if err != nil {
		s.Logger.Errorf("token获取失败: %s", err.Error())
		return
	}
	return
}
