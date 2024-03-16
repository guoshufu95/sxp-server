package api

import (
	"github.com/gin-gonic/gin"
	"sxp-server/app/service"
)

type RoleApi struct {
	Api
}

var rs = service.RoleService{}

func (a *RoleApi) CreateRole(c *gin.Context) {
	a.MakeApi(c)
	return
}
