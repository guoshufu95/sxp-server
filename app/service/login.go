package service

import (
	"sxp-server/app/dao"
	"sxp-server/app/model"
	"sxp-server/app/service/dto"
	"sxp-server/common/jwtToken"
	"sxp-server/common/utils"
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
	err, user := dao.GetUser(s.Db, req.Username)
	if err != nil {
		s.Logger.Error("根据用户名查询用户失败!")
		return
	}
	_, err = utils.CompareHashAndPassword(user.Password, req.Password)
	if err != nil {
		s.Logger.Error("用户密码不匹配")
		return
	}
	var role model.Role
	err = dao.GetRoleById(s.Db, user.RoleId, &role)
	if err != nil {
		s.Logger.Error("通过id查询用户信息失败")
		return
	}
	token, err = jwtToken.GenToken(req.Username, role.RoleKey, user.RoleId)
	if err != nil {
		s.Logger.Errorf("token获取失败: %s", err.Error())
		return
	}
	//更新登录时间
	err = dao.UpdateLoginTime(s.Db, user.ID)
	if err != nil {
		s.Logger.Error("更新登录时间失败")
		return
	}
	return
}
