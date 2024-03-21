package dto

import "sxp-server/app/model"

// CommonMenuReq
// @Description: insert || update
type CommonMenuReq struct {
	ParentId  int    `json:"parent_id"` //父级菜单id
	Name      string `json:"name"`      // 菜单名称
	Path      string `json:"path"`      // 路由地址
	Component string `json:"component"` // 组件路径
	Icon      string `json:"icon"`      // 图标
	OrderNum  int8   `json:"order_num"` // 排序
	Redirect  string `json:"redirect"`  // 重定向地址
	Hidden    int    `json:"is_hidden"` // 是否隐藏
}

// BuildData
//
//	@Description: 构建公共字段
//	@receiver c
//	@param data
func (c CommonMenuReq) BuildData(data *model.Menu) {
	data.ParentId = c.ParentId
	data.Name = c.Name
	data.Path = c.Path
	data.Component = c.Component
	data.Icon = c.Icon
	data.OrderNum = c.OrderNum
	data.Redirect = c.Redirect
	data.Hidden = c.Hidden
}

// CreateMenuReq
// @Description: 创建菜单入参
type CreateMenuReq struct {
	CommonMenuReq
}

// BuildCreateData
//
//	@Description: 构建入库数据
//	@receiver c
//	@param data
func (c CreateMenuReq) BuildCreateData(data *model.Menu) {
	c.CommonMenuReq.BuildData(data)
}

// UpdateMenuReq
// @Description: 更新菜单入参
type UpdateMenuReq struct {
	Id int `json:"id"` //菜单id
	CommonMenuReq
}

// BuildUpdateData
//
//	@Description: 构建更新数据入参
//	@receiver c
//	@param data
func (c UpdateMenuReq) BuildUpdateData(data *model.Menu) {
	c.CommonMenuReq.BuildData(data)
	data.ID = uint(c.Id)
}

// DeleteMenuReq
// @Description: 删除菜单入参，id
type DeleteMenuReq struct {
	Id int `json:"id"`
}
