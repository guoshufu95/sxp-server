package dto

import "sxp-server/app/model"

type CommonDeptReq struct {
	ParentId int    `json:"parentId"`
	Name     string `json:"name"`   //部门名称
	Leader   string `json:"leader"` //部门负责人
	Phone    string `json:"phone"`  //手机
	Email    string `json:"email"`  //邮箱
	Status   int    `json:"status"` //状态
}

// BuildData
//
//	@Description: 构建公共入库参数
//	@receiver c
//	@param data
func (c CommonDeptReq) BuildData(data *model.Dept) {
	data.ParentId = uint(c.ParentId)
	data.Name = c.Name
	data.Leader = c.Leader
	data.Phone = c.Phone
	data.Email = c.Email
	data.Status = c.Status
}

// CreateDeptReq
// @Description: 创建部门入参
type CreateDeptReq struct {
	CommonDeptReq
}

// BuildCreateData
//
//	@Description: 构造create入库参数
//	@receiver c
//	@param data
func (c CreateDeptReq) BuildCreateData(data *model.Dept) {
	c.CommonDeptReq.BuildData(data)
}

// UpdateDeptReq
// @Description: 更新部门入参
type UpdateDeptReq struct {
	Id int `json:"id"`
	CommonDeptReq
}

// BuildUpdateData
//
//	@Description: 构造更新入库参数
//	@receiver c
//	@param data
func (c UpdateDeptReq) BuildUpdateData(data *model.Dept) {
	c.CommonDeptReq.BuildData(data)
	data.ID = uint(c.Id)
}

// DeleteDeptReq
// @Description: 删除部门入参
type DeleteDeptReq struct {
	Id int `json:"id"`
}
