package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
	"sxp-server/app/dao"
	"sxp-server/app/model"
	"sxp-server/app/service/dto"
	ini "sxp-server/common/initial"
	"time"
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
func (s *UserService) ListUsers() (err error, res []dto.QueryRes) {
	err, users := dao.Users(s.Db)
	if err != nil {
		s.Logger.Error("查询用户列表失败")
		return
	}
	res = make([]dto.QueryRes, 0)
	dto.BuildQueryRes(&users, &res)
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

// GetUsersByParams
//
//	@Description: 条件查询
//	@receiver s
//	@param req
//	@return err
//	@return users
func (s *UserService) GetUsersByParams(req dto.QueryByParamsReq, users *[]model.User) (err error, res []dto.QueryRes) {
	db := s.buildQuery(s.Db, req)
	err = dao.GetUsersByParams(db, users)
	res = make([]dto.QueryRes, 0)
	dto.BuildQueryRes(users, &res)
	if err != nil {
		s.Logger.Error("条件查询user失败")
	}
	return
}

// buildQuery
//
//	@Description: 构建条件查询参数
//	@receiver s
//	@param db
//	@param req
//	@return *gorm.DB
func (s *UserService) buildQuery(db *gorm.DB, req dto.QueryByParamsReq) *gorm.DB {
	if req.UserName != "" {
		db = db.Where(fmt.Sprintf("username like \"%s\" or username like \"%s\" or username like \"%s\" or username = \"%s\"",
			"%"+req.UserName+"%",
			"%"+req.UserName,
			req.UserName+"%",
			req.UserName))
	}
	if req.Phone != "" {
		db = db.Where(fmt.Sprintf("phone like \"%s\" or phone like \"%s\" or phone like \"%s\" or phone = \"%s\"",
			"%"+req.Phone+"%",
			"%"+req.Phone,
			req.Phone+"%",
			req.Phone))
	}
	if req.Status != "" {
		if req.Status == "在线" {
			db = db.Where("status = ?", 1)
		} else {
			db = db.Where("status = ?", 0)
		}

	}
	return db
}

// GetUserById
//
//	@Description: 通过id返回用户信息
//	@receiver s
//	@param id
//	@return err
//	@return user
func (s *UserService) GetUserById(id int) (err error, res dto.QueryRes0) {
	var user model.User
	err = dao.GetUserById(s.Db, id, &user)
	res.Username = user.Username
	res.Id = user.ID
	res.Sex = user.Sex
	res.Email = user.Email
	res.NickName = user.NickName
	res.Phone = user.Phone
	res.LoginType = user.LoginType
	if user.LastLoginTime != nil {
		res.LastLoginTime = user.LastLoginTime.Format(time.DateTime)
	}

	res.Remark = user.Remark
	if user.Status == 1 {
		res.Status = "在线"
	} else {
		res.Status = "下线"
	}
	res.IsSuper = user.IsSuper
	if err != nil {
		s.Logger.Error("通过id查询用户信息失败")
	}
	dt := make([]uint, 0)
	for _, dept := range user.Depts {
		dt = append(dt, dept.ID)
	}
	res.Depts = dt
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
	if u.ID != 0 && u.ID != uint(req.Id) {
		err = errors.New("用户名重复，请重新设置")
		return
	}
	var depts []model.Dept
	err = dao.GetDeptsById(s.Db, req.DeptIds, &depts)
	if err != nil {
		s.Logger.Error("通过id列表查询部门信息失败")
		return
	}
	user.Depts = depts
	err = dao.ReplaceUserDept(s.Db, user)
	if err != nil {
		s.Logger.Error("更新user管理dept失败")
		return
	}
	err = dao.UpdateUser(db, user)
	if err != nil {
		s.Logger.Error("更新用户失败")
		return
	}
	return
}

// UpdateStatus
//
//	@Description: 更新用户在线状态
//	@receiver s
//	@param req
//	@return err
func (s *UserService) UpdateStatus(req dto.UpdateStatusReq) (err error) {
	err = dao.UpdateStatusById(s.Db, uint(req.Id), req.Status)
	if err != nil {
		s.Logger.Error("查询用户失败")
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
	var user model.User
	err = dao.GetUserById(s.Db, id, &user)
	if err != nil {
		s.Logger.Error("查询用户失败")
		return
	}
	err = dao.DeleteUerById(s.Db, user)
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
