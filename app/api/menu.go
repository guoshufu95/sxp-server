package api

import (
	"github.com/gin-gonic/gin"
	serv "sxp-server/app/service"
	"sxp-server/app/service/dto"
	"sxp-server/common/model"
)

type MenuApi struct {
	Api
}

var ms serv.MenuService

// GetMenus
//
//	@Description: 获取所有菜单
//	@receiver a
//	@param c
func (a *MenuApi) GetMenus(c *gin.Context) {
	a.BuildApi(c).BuildService(&ms.Service)
	err, menus := ms.ListMenu()
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	a.Response("成功获取所有菜单", menus)
}

// GetMenusByRole
//
//	@Description: 返回当前用户展示菜单
//	@receiver a
//	@param c
func (a *MenuApi) GetMenusByRole(c *gin.Context) {
	a.BuildApi(c).BuildService(&ms.Service)
	v := c.MustGet("sxp-claims")
	claims := v.(*model.MyClaims)
	err, menus := ms.GetRoleMenus(claims.RoleKey, claims.RoleId)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	a.Response("成功获取当前用户菜单", menus)
}

// CreateMenu
//
//	@Description: 创建菜单
//	@receiver a
//	@param c
func (a *MenuApi) CreateMenu(c *gin.Context) {
	a.BuildApi(c).BuildService(&ms.Service)
	var req dto.CreateMenuReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	err = ms.CreateMenu(req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	a.Response("创建菜单成功", "success")
}

// UpdateMenu
//
//	@Description: 更新菜单
//	@receiver a
//	@param c
func (a *MenuApi) UpdateMenu(c *gin.Context) {
	a.BuildApi(c).BuildService(&ms.Service)
	var req dto.UpdateMenuReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}

	err = ms.UpdateMenu(req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	a.Response("更新菜单成功", "success")
}

// DeleteMenu
//
//	@Description: 删除菜单
//	@receiver a
//	@param c
func (a *MenuApi) DeleteMenu(c *gin.Context) {
	a.BuildApi(c).BuildService(&ms.Service)
	var req dto.DeleteMenuReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	err = ms.DeleteMenu(req.Id)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	a.Response("删除菜单成功", "success")
}
