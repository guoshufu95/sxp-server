package service

import (
	"sxp-server/app/service/dto"
	"sxp-server/common/jwtToken"
)

type LoginService struct {
	Service
}

// Login
//
//	@Description: 登录校验相关
//	@receiver s
//	@param req
//	@return err
func (s *LoginService) Login(req dto.LoginReq) (err error, token string) {
	//err, user := login.GetUser(s.Db, req.Username)
	//if err != nil {
	//	s.Logger.Error("根据用户名查询用户失败!")
	//	return
	//}
	//_, err = utils.CompareHashAndPassword(user.PassWord, req.Password)
	//if err != nil {
	//	s.Logger.Error("用户密码不匹配")
	//	return
	//}
	token, err = jwtToken.GenToken(req.Username)
	if err != nil {
		s.Logger.Errorf("token获取失败: %s", err.Error())
		return
	}
	return
}
