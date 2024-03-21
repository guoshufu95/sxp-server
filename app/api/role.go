package api

import (
	"github.com/gin-gonic/gin"
	"sxp-server/app/service"
	"sxp-server/app/service/dto"
)

type RoleApi struct {
	Api
}

var rs = service.RoleService{}

// ListRoles
//
//	@Description: 角色列表
//	@receiver a
//	@param c
func (a *RoleApi) ListRoles(c *gin.Context) {
	a.BuildApi(c).BuildService(&rs.Service)
	err, roles := rs.ListRoles()
	if err != nil {
		a.Logger.Error(err)
		a.ResponseError(err)
		return
	}
	a.Response("success", roles)
}

// CreateRole
//
//	@Description: 创建角色
//	@receiver a
//	@param c
func (a *RoleApi) CreateRole(c *gin.Context) {
	a.BuildApi(c).BuildService(&rs.Service)
	var req dto.CreateRoleReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	err = rs.Auth(c)
	if err != nil {
		a.Logger.Error(err)
		a.ResponseError(err)
		return
	}
	err = rs.CreateRole(req)
	if err != nil {
		a.Logger.Error(err)
		a.ResponseError(err)
		return
	}
	a.Response("success", nil)
}

// UpdateRole
//
//	@Description:
//	@receiver a
//	@param c
func (a *RoleApi) UpdateRole(c *gin.Context) {
	a.BuildApi(c).BuildService(&rs.Service)
	var req dto.UpdateRoleReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	err = rs.Auth(c)
	if err != nil {
		a.Logger.Error(err)
		a.ResponseError(err)
		return
	}
	err = rs.UpdateRole(req)
	if err != nil {
		a.Logger.Error(err)
		a.ResponseError(err)
		return
	}
	a.Response("success", nil)
}

func (a *RoleApi) DeleteRole(c *gin.Context) {
	a.BuildApi(c).BuildService(&rs.Service)
}
