package dao

import (
	"gorm.io/gorm"
	"sxp-server/app/model"
)

// Users
//
//	@Description: 返回用户列表
//	@param db
//	@return err
//	@return users
func Users(db *gorm.DB) (err error, users []model.User) {
	err = db.Debug().Find(&users).Error
	return
}

// GetAuth
//
//	@Description: 根据roleId查询user信息
//	@param db
//	@param roleId
//	@return err
//	@return user
func GetAuth(db *gorm.DB, roleId int) (err error, user model.User) {
	db.Table("user").Where("role_id = ?", roleId).Debug().Find(&user)
	return
}

// GetUserByName
//
//	@Description: 根据用户名查询用户信息
//	@param db
//	@param name
//	@return err
//	@return user
func GetUserByName(db *gorm.DB, name string) (err error, user model.User) {
	err = db.Table("user").Where("username = ?", name).Find(&user).Error
	return
}

// GetUserById
//
//	@Description: 根据id查询用户信息
//	@param db
//	@param id
//	@return err
//	@return user
func GetUserById(db *gorm.DB, id int, user *model.User) (err error) {
	err = db.Find(&user, id).Error
	return
}

// CreateUser
//
//	@Description: 创建用户
//	@param db
//	@param user
//	@return err
func CreateUser(db *gorm.DB, user model.User) (err error) {
	err = db.Debug().Create(&user).Error
	return
}
