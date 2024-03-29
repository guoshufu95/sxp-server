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
	var (
		user  model.User
		depts []model.Dept
	)
	req.BuildCreateData(&user)
	err = dao.GetDeptsByIds(s.Db, req.DeptIds, &depts)
	if err != nil {
		s.Logger.Error("获取部门信息失败")
		return
	}
	user.Depts = depts
	err = dao.CreateUser(s.Db, user)
	if err != nil {
		s.Logger.Error("创建用户失败")
		return
	}
	return
}

// UpdateUser
//
//	@Description: 更新user
//	@receiver s
//	@param req
//	@return err
func (s *UserService) UpdateUser(req dto.UpdateUserReq) (err error) {
	var user model.User
	req.BuildUpdateData(&user)
	db := s.Db
	db.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	err, u := dao.GetUserByName(db, req.Username)
	if err != nil {
		s.Logger.Error("通过name查询失败")
		return
	}
	if u.ID != uint(req.Id) {
		err = errors.New("用户名重复，请重新设置")
		return
	}
	err = dao.UpdateUser(db, user)
	if err != nil {
		s.Logger.Error("更新用户失败")
		return
	}
	return
}

// DeleteUser
//
//	@Description: 删除用户
//	@receiver s
//	@param id
//	@return err
func (s *UserService) DeleteUser(id int) (err error) {
	err = dao.DeleteUerById(s.Db, id)
	if err != nil {
		s.Logger.Error("删除用户失败")
		return
	}
	return
}

// CasbinPermission
//
//	@Description: 测试用
//	@receiver s
//	@param roleId
//	@return err
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
	// 加入casbin表
	_, err = e.AddNamedPolicies("p", policys)
	return
}
