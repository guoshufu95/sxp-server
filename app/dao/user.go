package dao

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sxp-server/app/model"
)

// Users
//
//	@Description: 返回用户列表
//	@param db
//	@return err
//	@return users
func Users(db *gorm.DB) (err error, users []model.User) {
	err = db.Debug().Preload("Depts").Find(&users).Error
	return
}

// GetAuth
//
//	@Description: 根据roleId查询user信息
//	@param db
//	@param roleId
//	@return err
//	@return user
func GetAuth(db *gorm.DB, name string) (err error, user model.User) {
	db.Table("user").Where("username = ?", name).Debug().Find(&user)
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
	err = db.Table("user").Debug().Where("username = ?", name).Find(&user).Error
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
	err = db.Debug().Preload("Depts").Find(&user, id).Error
	return
}

// GetUsersByParams
//
//	@Description: 条件查询
//	@param db
//	@param req
//	@param user
//	@return err
func GetUsersByParams(db *gorm.DB, user *[]model.User) (err error) {
	err = db.Model(&model.User{}).Debug().Find(&user).Error
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

// UpdateUser
//
//	@Description: 更新user
//	@param db
//	@param user
//	@return err
func UpdateUser(db *gorm.DB, user model.User) (err error) {
	err = db.Debug().Model(&model.User{}).Where("id = ?", user.ID).Updates(&user).Error
	return
}

// DeleteUerById
//
//	@Description: 通过id删除用户
//	@param db
//	@param id
//	@return err
func DeleteUerById(db *gorm.DB, user model.User) (err error) {
	err = db.Debug().Unscoped().Select(clause.Associations).Delete(&user).Error
	return
}
