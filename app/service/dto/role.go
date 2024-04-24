package dto

import (
	"strconv"
	"sxp-server/app/model"
)

// RoleCommonReq
// @Description: 请求的公共参数
type RoleCommonReq struct {
	Name     string `json:"name"`
	RoleKey  string `json:"roleKey"`
	Label    string `json:"label"`
	Status   int    `json:"status"`
	RoleSort string `json:"roleSort"`
	MenuIds  []int  `json:"menuIds"`
	DeptIds  []int  `json:"deptIds"`
}

func (c RoleCommonReq) BuildData(data *model.Role) {
	data.Name = c.Name
	data.RoleKey = c.RoleKey
	data.Label = c.Label
	data.Status = c.Status
	rs, _ := strconv.Atoi(c.RoleSort)
	data.RoleSort = rs
}

// CreateRoleReq
// @Description: 新增角色请求参数
type CreateRoleReq struct {
	RoleCommonReq
}

// BuildCreateData
//
//	@Description: 构造入库参数
//	@receiver c
//	@param data
func (c CreateRoleReq) BuildCreateData(data *model.Role) {
	c.RoleCommonReq.BuildData(data)
}

// UpdateRoleReq
// @Description: 更新role入参
type UpdateRoleReq struct {
	Id int `json:"id"`
	RoleCommonReq
}

type UpdateRoleStatusReq struct {
	Id     int `json:"id"`
	Status int `json:"status"`
}

// BuildUpdateData
//
//	@Description: 构造更新参数
//	@receiver c
//	@param data
func (c UpdateRoleReq) BuildUpdateData(data *model.Role) {
	c.RoleCommonReq.BuildData(data)
	data.ID = uint(c.Id)
}

// DeleteRoleReq
// @Description:删除role入参
type DeleteRoleReq struct {
	Id int `json:"id"`
}

// QueryRoleByParams
// @Description: 角色列表条件查询
type QueryRoleByParams struct {
	Name    string `json:"name"`
	RoleKey string `json:"roleKey"`
	Status  string `json:"status"`
}

// GetRoleByIdReq
// @Description: id查询入参
type GetRoleByIdReq struct {
	Id int `json:"id"`
}

// GetRoleByIdRes
// @Description:id查询返回
type GetRoleByIdRes struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	RoleKey  string `json:"roleKey"`
	Label    string `json:"label"`
	Status   int    `json:"status"`
	RoleSort string `json:"roleSort"`
	MenuIds  []uint `json:"menuIds"`
	DeptIds  []uint `json:"deptIds"`
}
