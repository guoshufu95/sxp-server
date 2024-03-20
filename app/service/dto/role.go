package dto

import "sxp-server/app/model"

// RoleCommonReq
// @Description: 请求的公共参数
type RoleCommonReq struct {
	Name     string `json:"name"`
	RoleKey  string `json:"roleKey"`
	Label    string `json:"label"`
	Status   int    `json:"status"`
	RoleSort int    `json:"roleSort"`
	MenuIds  []int  `json:"menus"`
	DeptIds  []int  `json:"deptIds"`
}

func (c RoleCommonReq) BuildData(data *model.Role) {
	data.Name = c.Name
	data.RoleKey = c.RoleKey
	data.Label = c.Label
	data.Status = c.Status
	data.RoleSort = c.RoleSort
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

func (c UpdateRoleReq) BuildUpdateData(data *model.Role) {
	c.RoleCommonReq.BuildData(data)
	data.ID = uint(c.Id)
}
