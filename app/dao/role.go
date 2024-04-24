package dao

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	err = db.Table("role").Debug().Preload("Menus").Preload("Depts").Where("id = ?", id).Find(&role).Error
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

// UpdateRoleDepts
//
//	@Description: 更新role关联的部门
//	@param db
//	@param role
//	@return err
func UpdateRoleDepts(db *gorm.DB, role model.Role) (err error) {
	err = db.Debug().Model(&role).Association("Depts").Replace(role.Depts)
	return
}

// UpdateRoleMenus
//
//	@Description: 更新role关联的菜单
//	@param db
//	@param role
//	@return err
func UpdateRoleMenus(db *gorm.DB, role model.Role) (err error) {
	err = db.Debug().Model(&role).Association("Menus").Replace(role.Menus)
	return
}

// DeleteRoleMenus
//
//	@Description: 删除role绑定的菜单
//	@param db
//	@param data
//	@return err
func DeleteRoleMenus(db *gorm.DB, data model.Role) (err error) {
	err = db.Debug().Model(&model.Role{}).Where("id = ?", data.ID).Delete(&data.Menus).Error
	return
}

// DeleteRoleDepts
//
//	@Description: 删除role绑定的部门
//	@param db
//	@param data
//	@return err
func DeleteRoleDepts(db *gorm.DB, data model.Role) (err error) {
	err = db.Debug().Model(&model.Role{}).Where("id = ?", data.ID).Delete(&data.Depts).Error
	return
}

// UpdateRole
//
//	@Description: 更新role
//	@param db
//	@param data
//	@return err
func UpdateRole(db *gorm.DB, data model.Role) (err error) {
	err = db.Debug().Model(&data).Where("id = ?", data.ID).Updates(&data).Error
	if err != nil {
		return
	}
	err = db.Debug().Model(&data).Association("Menus").Replace(data.Menus)
	if err != nil {
		return
	}
	err = db.Debug().Model(&data).Association("Depts").Replace(data.Depts)
	if err != nil {
		return
	}
	return
}

// DeleteRoleById
//
//	@Description: 删除role
//	@param db
//	@param id
//	@return err
func DeleteRoleById(db *gorm.DB, role model.Role) (err error) {
	err = db.Debug().Select(clause.Associations).Unscoped().Delete(&role).Error
	return
}

// GetRoleByDepts
//
//	@Description: 通过detps关联查询出roles
//	@param db
//	@param depts[]
//	@param roles
//	@return err
func GetRoleByDepts(db *gorm.DB, dept model.Dept, roles *[]model.Role) (err error) {
	err = db.Model(&dept).Debug().Association("Roles").Find(&roles)
	return
}

// GetRoleByParams
//
//	@Description: 条件查询
//	@param db
//	@param roles
//	@return err
func GetRoleByParams(db *gorm.DB, roles *[]model.Role) (err error) {
	err = db.Debug().Find(&roles).Error
	return
}

// UpdateRoleStatus
//
//	@Description: 更新启用状态
//	@param db
//	@param id
//	@param status
//	@return err
func UpdateRoleStatus(db *gorm.DB, id, status int) (err error) {
	err = db.Debug().Model(&model.Role{}).Where("id = ?", id).Update("status", status).Error
	return
}
