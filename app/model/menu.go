package model

import "gorm.io/gorm"

// Menu
// @Description:数据库字段
type Menu struct {
	gorm.Model
	ParentId int    `json:"parentId" gorm:"type:int(4);comment:父级菜单id;"`
	Name     string `json:"label" gorm:"type:varchar(100);comment:菜单名;unique;not null"` // 菜单名称
	Path     string `json:"path" gorm:"type:varchar(255);comment:菜单路由地址;"`              // 路由地址
	//Component  string  `json:"component" gorm:"type:varchar(100);comment:组件路径;"`           // 组件路径
	Icon       string  `json:"icon" gorm:"type:varchar(100);comment:图标;"`        // 图标
	OrderNum   string  `json:"orderNum" gorm:"int(4);comment:排序;"`               // 排序
	Redirect   string  `json:"redirect" gorm:"type:varchar(255);comment:重定向地址;"` // 重定向地址
	Hidden     int     `json:"hidden" gorm:"type:int(2);comment:是否隐藏;"`          // 是否隐藏
	Children   []Menu  `json:"children" gorm:"-"`                                //子菜单
	Roles      []*Role `json:"roles" gorm:"many2many:role_menu;"`                //关联角色表
	CreateTime string  `json:"createTime" gorm:"-"`
}

func (Menu) TableName() string {
	return "menu"
}
