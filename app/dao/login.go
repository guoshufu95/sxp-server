package dao

import (
	"gorm.io/gorm"
	"sxp-server/app/model"
	"time"
)

// GetUser
//
//	@Description: 根据用户名查询
//	@param db
//	@param name
//	@return err
func GetUser(db *gorm.DB, name string) (err error, user model.User) {
	err = db.Table("user").Where("username = ?", name).Preload("Depts").Find(&user).Error
	return
}

// UpdateLoginTime
//
//	@Description: 更新登录时间
//	@param db
//	@param userId
//	@return err
func UpdateLoginTime(db *gorm.DB, userId uint) (err error) {
	err = db.Table("user").Debug().Where("id = ?", userId).Update("last_login_time", time.Now()).Error
	return
}
