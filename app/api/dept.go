package api

import (
	"github.com/gin-gonic/gin"
	serv "sxp-server/app/service"
	"sxp-server/app/service/dto"
)

type DeptApi struct {
	Api
}

var ds serv.DeptService

// GetDepts
//
//	@Description: 返回部门列表
//	@receiver a
//	@param c
func (a *DeptApi) GetDepts(c *gin.Context) {
	a.BuildApi(c).BuildService(&ds.Service)
	err, dept := ds.GetDept()
	if err != nil {
		a.Logger.Error(err)
		a.ResponseError(err)
		return
	}
	a.Response("成功返回部门列表!", dept)
}

// InsertDept
//
//	@Description: 创建部门
//	@receiver a
//	@param c
func (a *DeptApi) InsertDept(c *gin.Context) {
	a.BuildApi(c).BuildService(&ds.Service)
	var req dto.CreateDeptReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	err = ds.Auth(c) //只有admin才能创建
	if err != nil {
		a.Logger.Error(err)
		a.ResponseError(err)
		return
	}
	err = ds.CreateDept(req)
	if err != nil {
		a.Logger.Error(err)
		a.ResponseError(err)
		return
	}
	a.Response("成功创建部门!", nil)
}

func (a *DeptApi) UpdateDept(c *gin.Context) {
	a.BuildApi(c).BuildService(&ds.Service)
	var req dto.UpdateDeptReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	err = ds.Auth(c) //只有admin才能更新
	if err != nil {
		a.Logger.Error(err)
		a.ResponseError(err)
		return
	}
	err = ds.UpdateDept(req)
	if err != nil {
		a.Logger.Error(err)
		a.ResponseError(err)
		return
	}
	a.Response("成功更新部门!", nil)
}

// DeleteDept
//
//	@Description: 删除部门
//	@receiver a
//	@param c
func (a *DeptApi) DeleteDept(c *gin.Context) {
	a.BuildApi(c).BuildService(&ds.Service)
	var req dto.UpdateDeptReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.Logger.Error(err.Error())
		a.ResponseError(err)
		return
	}
	err = ds.Auth(c) //只有admin才能删除
	if err != nil {
		a.Logger.Error(err)
		a.ResponseError(err)
		return
	}
	err = ds.DeleteDept(req.Id)
	if err != nil {
		a.Logger.Error(err)
		a.ResponseError(err)
		return
	}
	a.Response("删除成功!", nil)
}
