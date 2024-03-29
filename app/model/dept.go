package model

import "gorm.io/gorm"

// Dept
// @Description: 部门(用户组)
type Dept struct {
	gorm.Model
	ParentId uint   `json:"parentId" gorm:"type:int(5);comment:父id"`       //上级部门id
	Name     string `json:"name" gorm:"type:varchar(58);comment:部门名称"`     //部门名称
	Leader   string `json:"leader" gorm:"type:varchar(100);comment:部门负责人"` //部门负责人
	Phone    string `json:"phone" gorm:"type:varchar(15);comment:电话"`      //手机
	Email    string `json:"email" gorm:"type:varchar(100);comment:邮箱"`     //邮箱
	Status   int    `json:"status" gorm:"type:varchar(2);comment:启用状态"`    //状态
	Roles    []Role `json:"roles" gorm:"many2many:role_dept;"`
	Users    []User `json:"users" gorm:"many2many:user_dept;"`
	Children []Dept `json:"children" gorm:"-"`
}

func (Dept) TableName() string {
	return "dept"
}
