package model

import (
	"gorm.io/gorm"
	"time"
)

// User
// @Description: user表字段
type User struct {
	gorm.Model
	Username      string     `json:"username" gorm:"type:varchar(255);comment:用户名;unique;not null"`
	Password      string     `json:"password" gorm:"type:varchar(100);comment:密码;not null;"`
	NickName      string     `json:"nick_name" gorm:"type:varchar(255);comment:昵称;unique"`
	Sex           string     `json:"sex" gorm:"type:varchar(10);comment:性别;"`
	Email         string     `json:"email" gorm:"type:varchar(100);comment:邮箱;"`
	Phone         string     `json:"phone" gorm:"type:varchar(50);comment:电话;"`
	LoginType     int        `json:"login_type" gorm:"type:tinyint(4);comment:登录类型;"`
	LastLoginTime *time.Time `json:"last_login_time" gorm:"comment:上次登录时间"`
	Remark        string     `json:"remark" gorm:"type:varchar(100);comment:描述;"`
	Status        string     `json:"status" gorm:"type:int(4);comment:启用状态;"`
	IsSuper       int        `json:"is_super" gorm:"int(2);comment:是否是超级管理员;"`
	RoleId        string     `json:"roleId" gorm:"int(2);comment:用户所属的角色id"`
}

func (User) TableName() string {
	return "user"
}
