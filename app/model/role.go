package model

import "gorm.io/gorm"

// Role
// @Description: role数据库字段
type Role struct {
	gorm.Model
	Name     string `json:"name" gorm:"type:varchar(100);comment:角色名;unique;not null"`
	Label    string `json:"label" gorm:"type:varchar(100);comment:标签;"`
	Status   int    `json:"is_disable" gorm:"type:int(2);comment:启用状态;"`
	RoleSort int    `json:"roleSort" gorm:"type:int(4);comment:排序字段;"`
	Menus    []Menu `json:"menus" gorm:"many2many:role_menu;"`
	Users    []User `json:"users" gorm:"many2many:user_role;"`
}

func (Role) TableName() string {
	return "role"
}
