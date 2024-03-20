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

// GetMenusByIds
//
//	@Description: 通过ids返回列表
//	@param db
//	@param ids[]int
//	@param menus
//	@return err
func GetMenusByIds(db *gorm.DB, ids []int, menus *[]model.Menu) (err error) {
	err = db.Debug().Find(&menus, ids).Error
	return
}
