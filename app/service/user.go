package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"sxp-server/app/dao"
	"sxp-server/app/model"
	"sxp-server/app/service/dto"
	ini "sxp-server/common/initial"
	cm "sxp-server/common/model"
)

type UserService struct {
	Service
}

// ListUsers
//
//	@Description: 用户列表
//	@receiver s
//	@return err
//	@return users
func (s *UserService) ListUsers() (err error, users []model.User) {
	err, users = dao.Users(s.Db)
	if err != nil {
		s.Logger.Error("查询用户列表失败")
		return
	}
	return
}

// Auth
//
//	@Description:
//	@param db
//	@param c
//	@return err
//	@return flag
func (s *UserService) Auth(c *gin.Context) (err error) {
	v, ok := c.Get("sxp-claims")
	if !ok {
		err = errors.New("无法获取claims")
		return
	}
	claims := v.(*cm.MyClaims)
	err, user := dao.GetAuth(s.Db, claims.RoleId)
	if err != nil {
		err = errors.New("获取当前登录用户信息失败")
		return
	}
	if user.IsSuper == 0 {
		err = errors.New("权限不足，只有超级管理员才能创建用户")
		return
	}
	return
}

// GetUserByName
//
//	@Description: 通过username返回用户信息
//	@receiver s
//	@param name
//	@return err
func (s *UserService) GetUserByName(name string) (err error) {
	err, user := dao.GetUserByName(s.Db, name)
	if err != nil {
		s.Logger.Error("根据name查询user失败")
		return
	}

	if user.ID != 0 {
		err = errors.New("用户名已存在")
		return
	}
	return
}

// GetUserById
//
//	@Description: 通过id返回用户信息
//	@receiver s
//	@param id
//	@return err
//	@return user
func (s *UserService) GetUserById(id int) (err error, user model.User) {
	err = dao.GetUserById(s.Db, id, &user)
	if err != nil {
		s.Logger.Error("通过id查询用户信息失败")
	}
	return
}

// CreateUser
//
//	@Description: 创建用户
//	@receiver s
//	@param req
//	@return err
func (s *UserService) CreateUser(req dto.CreateUserReq) (err error) {
	var user model.User
	req.BuildData(&user)
	err = dao.CreateUser(s.Db, user)
	if err != nil {
		s.Logger.Error("创建用户失败")
		return
	}
	return
}

func (s *UserService) CasbinPermission(roleId int) (err error) {
	var role model.Role
	err = dao.GetRoleById(s.Db, roleId, &role)
	if err != nil {
		s.Logger.Error("通过id获取角色信息失败")
		return
	}
	data := make([]gin.RouteInfo, 0)
	routes := ini.App.Engine.Routes()
	for _, v := range routes {
		if strings.Contains(v.Path, "task") { //测试用
			data = append(data, v)
		}
	}
	fmt.Print(data)
	policys := make([][]string, 0)
	e := ini.App.GetCasbin()
	for _, d := range data {
		policys = append(policys, []string{role.RoleKey, d.Path, d.Method})
	}
	_, err = e.AddNamedPolicies("p", policys)
	return
}
