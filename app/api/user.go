package api

import (
	"github.com/gin-gonic/gin"
	"sxp-server/app/service"
	"sxp-server/app/service/dto"
)

type UserApi struct {
	Api
}

var us service.UserService

// ListUsers
//
//	@Description: 获取用户列表
//	@receiver a
//	@param c
func (a UserApi) ListUsers(c *gin.Context) {
	a.BuildApi(c).BuildService(&us.Service)
	err, users := us.ListUsers()
	if err != nil {
		a.ResponseError(err)
		return
	}
	a.Response("返回用户列表成功", users)
}

// CreateUser
//
//	@Description: 创建user
//	@receiver a
//	@param c
func (a UserApi) CreateUser(c *gin.Context) {
	a.BuildApi(c).BuildService(&us.Service)
	var req = dto.CreateUserReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.ResponseError(err)
		return
	}
	//权限校验
	err = us.Auth(c)
	if err != nil {
		a.ResponseError(err)
		return
	}
	// 用户名校验
	err = us.GetUserByName(req.Username)
	if err != nil {
		a.ResponseError(err)
		return
	}
	// 数据库创建
	err = us.CreateUser(req)
	if err != nil {
		a.ResponseError(err)
		return
	}
	//todo casbin授权
	//err = us.CasbinPermission(req.RoleId)
	a.Response("创建用户成功", nil)
}

// GetById
//
//	@Description: 通过id查询用户信息
//	@receiver a
//	@param c
func (a UserApi) GetById(c *gin.Context) {
	a.BuildApi(c).BuildService(&us.Service)
	var req dto.GetUserByIdRequest
	err, user := us.GetUserById(req.Id)
	if err != nil {
		a.ResponseError(err)
		return
	}
	a.Response("获取用户信息成功", user)
}

// UpdateUser
//
//	@Description: 更新
//	@receiver a
//	@param c
func (a UserApi) UpdateUser(c *gin.Context) {
	a.BuildApi(c).BuildService(&us.Service)
	var req dto.UpdateUserReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.ResponseError(err)
		return
	}
	err = us.Auth(c) //权限
	if err != nil {
		a.ResponseError(err)
		return
	}
	err = us.UpdateUser(req)
	if err != nil {
		a.ResponseError(err)
		return
	}
	a.Response("更新数据成功", nil)
}

// DeleteUser
//
//	@Description: 删除
//	@receiver a
//	@param c
func (a UserApi) DeleteUser(c *gin.Context) {
	a.BuildApi(c).BuildService(&us.Service)
	var req dto.DeleteUserReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.ResponseError(err)
		return
	}
	err = us.Auth(c) //权限
	if err != nil {
		a.ResponseError(err)
		return
	}
	err = us.DeleteUser(req.Id)
	if err != nil {
		a.ResponseError(err)
		return
	}
	a.Response("用户删除成功", nil)
}
