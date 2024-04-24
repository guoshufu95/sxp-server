package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
func (a RoleApi) ListRoles(c *gin.Context) {
	a.BuildApi(c).BuildService(&rs.Service)
	err, roles := rs.ListRoles()
	if err != nil {
		a.ResponseError(http.StatusBadRequest, err)
		return
	}
	a.Response("success", roles)
}

// GetRoleById
//
//	@Description: 通过id查询返回详情
//	@receiver a
//	@param c
func (a RoleApi) GetRoleById(c *gin.Context) {
	a.BuildApi(c).BuildService(&rs.Service)
	var req dto.GetRoleByIdReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.ResponseError(http.StatusBadRequest, err)
		return
	}
	err, role := rs.GetRoleById(req.Id)
	if err != nil {
		a.ResponseError(http.StatusInternalServerError, err)
		return
	}
	a.Response("查询成功", role)

}

// GetRoleByParams
//
//	@Description: 条件查询
//	@receiver a
//	@param c
func (a RoleApi) GetRoleByParams(c *gin.Context) {
	a.BuildApi(c).BuildService(&rs.Service)
	var req dto.QueryRoleByParams
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.ResponseError(http.StatusBadRequest, err)
		return
	}
	err, roles := rs.RoleByParams(req)
	if err != nil {
		a.ResponseError(http.StatusInternalServerError, err)
		return
	}
	a.Response("查询成功", roles)
}

// CreateRole
//
//	@Description: 创建角色
//	@receiver a
//	@param c
func (a RoleApi) CreateRole(c *gin.Context) {
	a.BuildApi(c).BuildService(&rs.Service)
	var req dto.CreateRoleReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.ResponseError(http.StatusBadRequest, err)
		return
	}
	err = rs.Auth(c)
	if err != nil {
		a.ResponseError(http.StatusForbidden, err)
		return
	}
	err = rs.CreateRole(req)
	if err != nil {
		a.ResponseError(http.StatusInternalServerError, err)
		return
	}
	a.Response("创建角色成功", nil)
}

// UpdateRole
//
//	@Description:
//	@receiver a
//	@param c
func (a RoleApi) UpdateRole(c *gin.Context) {
	a.BuildApi(c).BuildService(&rs.Service)
	var req dto.UpdateRoleReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.ResponseError(http.StatusBadRequest, err)
		return
	}
	err = rs.Auth(c) //权限
	if err != nil {
		a.ResponseError(http.StatusForbidden, err)
		return
	}
	err = rs.UpdateRole(req)
	if err != nil {
		a.ResponseError(http.StatusInternalServerError, err)
		return
	}
	a.Response("更新角色成功", nil)
}

// UpdateStatus
//
//	@Description: 更新启用状态
//	@receiver a
//	@param c
func (a RoleApi) UpdateStatus(c *gin.Context) {
	a.BuildApi(c).BuildService(&rs.Service)
	var req dto.UpdateRoleStatusReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.ResponseError(http.StatusBadRequest, err)
		return
	}
	err = rs.UpdateStatus(req)
	if err != nil {
		a.ResponseError(http.StatusInternalServerError, err)
		return
	}
	a.Response("状态更新成功", nil)
}

// DeleteRole
//
//	@Description: 删除用户
//	@receiver a
//	@param c
func (a RoleApi) DeleteRole(c *gin.Context) {
	a.BuildApi(c).BuildService(&rs.Service)
	var req dto.UpdateRoleReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.ResponseError(http.StatusBadRequest, err)
		return
	}
	err = rs.Auth(c) //权限检查
	if err != nil {
		a.ResponseError(http.StatusForbidden, err)
		return
	}
	err = rs.DeleteRole(req.Id)
	if err != nil {
		a.ResponseError(http.StatusInternalServerError, err)
		return
	}
	a.Response("删除角色成功", nil)
}
