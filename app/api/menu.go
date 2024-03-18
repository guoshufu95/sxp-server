package api

import (
	"github.com/gin-gonic/gin"
	serv "sxp-server/app/service"
	"sxp-server/app/service/dto"
)

type MenuApi struct {
	Api
}

var ms serv.MenuService

// GetMenu
//
//	@Description: 获取菜单树
//	@receiver a
//	@param c
func (a *MenuApi) GetMenu(c *gin.Context) {

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

	return
}
