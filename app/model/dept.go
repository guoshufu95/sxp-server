package model

import "gorm.io/gorm"

// Dept
// @Description: 部门
type Dept struct {
	gorm.Model
	ParentId uint   `json:"parent_id"`                                     //上级部门id
	Name     string `json:"name" gorm:"type:varchar(58);comment:部门名称"`     //部门名称
	Leader   string `json:"leader" gorm:"type:varchar(100);comment:部门负责人"` //部门负责人
	Phone    string `json:"phone" gorm:"type:varchar(15);comment:电话"`      //手机
	Email    string `json:"email" gorm:"type:varchar(100);comment:邮箱"`     //邮箱
	Status   int    `json:"status" gorm:"type:varchar(2);comment:启用状态"`    //状态
	Roles    []Role `json:"roles" gorm:"many2many:role_dept;"`
	Children []Dept `json:"children"` //下级部门
}

func (Dept) TableName() string {
	return "dept"
}
