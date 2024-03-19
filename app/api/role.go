package api

import (
	"github.com/gin-gonic/gin"
	"sxp-server/app/service"
)

type RoleApi struct {
	Api
}

var rs = service.RoleService{}

// CreateRole
//
//	@Description: 创建角色
//	@receiver a
//	@param c
func (a *RoleApi) CreateRole(c *gin.Context) {
	a.BuildApi(c).BuildService(&rs.Service)
	return
}
