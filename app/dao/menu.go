package dao

import (
	"gorm.io/gorm"
	"sxp-server/app/model"
)

// ListMenus
//
//	@Description: 返回所有菜单
//	@param db
//	@param menus
//	@return err
func ListMenus(db *gorm.DB, menus *[]model.Menu) (err error) {
	err = db.Debug().Find(&menus).Error
	return
}

// RoleMenus
//
//	@Description: 返回当前用户的
//	@param db
//	@param id
//	@param menus
//	@return err
func RoleMenus(db *gorm.DB, id int, role *model.Role) (err error) {
	err = db.Model(&role).Where("id = ?", id).Preload("Menus").Find(&role).Error
	return
}
