package dto

import (
	"strconv"
	"sxp-server/app/model"
	"time"
)

// CommonMenuReq
// @Description: insert || update
type CommonMenuReq struct {
	ParentId  int    `json:"parentId"`  //父级菜单id
	Name      string `json:"label"`     // 菜单名称
	Path      string `json:"path"`      // 路由地址
	Component string `json:"component"` // 组件路径
	Icon      string `json:"icon"`      // 图标
	OrderNum  string `json:"orderNum"`  // 排序
	Redirect  string `json:"redirect"`  // 重定向地址
	Hidden    string `json:"hidden"`    // 是否隐藏
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
	//data.Component = c.Component
	data.Icon = c.Icon
	data.OrderNum = c.OrderNum
	data.Redirect = c.Redirect
	hidden, _ := strconv.Atoi(c.Hidden)
	data.Hidden = hidden
}

// CreateMenuReq
// @Description: 创建菜单入参
type CreateMenuReq struct {
	CommonMenuReq
}

// QueryMenusByParamReq
// @Description: 条件查询参数
type QueryMenusByParamReq struct {
	Name string `json:"name"`
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

// ListMenusRes
// @Description: menu返回数据
type ListMenusRes struct {
	Id       uint   `json:"id"`
	ParentId int    `json:"parentId" `
	Name     string `json:"label"`
	Path     string `json:"path"`
	//Component  string         `json:"component"`
	Icon       string         `json:"icon"`
	OrderNum   string         `json:"orderNum"`
	Redirect   string         `json:"redirect"`
	Hidden     int            `json:"hidden"`
	CreateTime string         `json:"createTime"`
	Children   []ListMenusRes `json:"children"`
}

// BuildResponse
//
//	@Description: 构造返回参数
//	@receiver l
//	@param menu
func (l *ListMenusRes) BuildResponse(menu model.Menu) {
	l.Id = menu.ID
	l.ParentId = menu.ParentId
	l.Name = menu.Name
	l.Path = menu.Path
	//l.Component = menu.Component
	l.Icon = menu.Icon
	l.OrderNum = menu.OrderNum
	l.Redirect = menu.Redirect
	l.Hidden = menu.Hidden
	l.CreateTime = menu.CreatedAt.Format(time.DateTime)
}

// GetMenuByIdReq
// @Description: id参数
type GetMenuByIdReq struct {
	Id int `json:"id"`
}
