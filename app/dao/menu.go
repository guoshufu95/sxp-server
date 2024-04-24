package dao

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
func RoleMenus(db *gorm.DB, ids []int, roles *[]model.Role) (err error) {
	err = db.Preload("Menus").Find(&roles, ids).Error
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

// CreateMenu
//
//	@Description: 创建菜单
//	@param db
//	@param menu
//	@return err
func CreateMenu(db *gorm.DB, menu model.Menu) (err error) {
	err = db.Debug().Create(&menu).Error
	return
}

// QueryMenuByParam
//
//	@Description: 返回条件查询数据
//	@param db
//	@param menu
//	@return err
func QueryMenuByParam(db *gorm.DB, menus *[]model.Menu) (err error) {
	err = db.Debug().Find(&menus).Error
	return
}

// UpdateMenu
//
//	@Description: 更新菜单
//	@param db
//	@param menu
//	@return err
func UpdateMenu(db *gorm.DB, id uint, m map[string]interface{}) (err error) {
	err = db.Debug().Model(&model.Menu{}).Where("id = ?", id).Updates(&m).Error
	return
}

// GetMenusById
//
//	@Description: 通过id查
//	@param db
//	@param id
//	@param menu
//	@return err
func GetMenusById(db *gorm.DB, id uint, menu *model.Menu) (err error) {
	err = db.Debug().Find(&menu, id).Error
	return
}

// DeleteMenuByIds
//
//	@Description: 删除菜单
//	@param db
//	@param ids
//	@return err
func DeleteMenuByIds(db *gorm.DB, menu model.Menu) (err error) {
	err = db.Debug().Select(clause.Associations).Unscoped().Delete(&menu).Error
	return
}
