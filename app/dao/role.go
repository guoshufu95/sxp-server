package dao

import (
	"gorm.io/gorm"
	"sxp-server/app/model"
)

// GetRoleById
//
//	@Description: 通过id获取
//	@param db
//	@param id
//	@param role
//	@return err
func GetRoleById(db *gorm.DB, id int, role *model.Role) (err error) {
	err = db.Table("role").Debug().Where("id = ?", id).Find(&role).Error
	return
}

// ListRoles
//
//	@Description: 角色列表
//	@param db
//	@param roles
//	@return err
func ListRoles(db *gorm.DB, roles *[]model.Role) (err error) {
	err = db.Debug().Preload("Menus").Preload("Depts").Find(&roles).Error
	return
}

// CreateRole
//
//	@Description: 创建role
//	@param db
//	@param data
//	@return err
func CreateRole(db *gorm.DB, data model.Role) (err error) {
	err = db.Debug().Create(&data).Error
	return
}

// DeleteRoleMenus
//
//	@Description: 删除role绑定的菜单
//	@param db
//	@param data
//	@return err
func DeleteRoleMenus(db *gorm.DB, data model.Role) (err error) {
	err = db.Debug().Association("Menus").Delete(data.Menus)
	return
}

// DeleteRoleDepts
//
//	@Description: 删除role绑定的部门
//	@param db
//	@param data
//	@return err
func DeleteRoleDepts(db *gorm.DB, data model.Role) (err error) {
	err = db.Debug().Association("Depts").Delete(data.Depts)
	return
}

// UpdateRole
//
//	@Description: 更新role
//	@param db
//	@param data
//	@return err
func UpdateRole(db *gorm.DB, data model.Role) (err error) {
	err = db.Model(&model.Role{}).Debug().Updates(&data).Error
	return
}
