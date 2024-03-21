package dao

import (
	"gorm.io/gorm"
	"sxp-server/app/model"
)

// GetDeptsByIds
//
//	@Description: 通过ids返回部门列表
//	@param db
//	@param ids
//	@param depts
//	@return err
func GetDeptsByIds(db *gorm.DB, ids []int, depts *[]model.Dept) (err error) {
	err = db.Debug().Find(&depts, ids).Error
	return
}

// GetDeptById
//
//	@Description: 通过id查找部门
//	@param db
//	@param ids
//	@param depts
//	@return err
func GetDeptById(db *gorm.DB, id uint, dept *model.Dept) (err error) {
	err = db.Debug().Model(&model.Dept{}).Where("id = ?", id).Find(&dept).Error
	return
}

// GetDeptByName
//
//	@Description: 通过部门名查询
//	@param db
//	@param name
//	@return err
//	@return dept
func GetDeptByName(db *gorm.DB, name string) (err error, dept model.Dept) {
	err = db.Debug().Model(&model.Dept{}).Where("name = ?", name).Find(&dept).Error
	return
}

// GetAllDepts
//
//	@Description: 返回所有部门
//	@param db
//	@param depts
//	@return err
func GetAllDepts(db *gorm.DB, depts *[]model.Dept) (err error) {
	err = db.Debug().Find(&depts).Error
	return
}

// CreateDept
//
//	@Description: 创建部门
//	@param db
//	@param data
//	@return err
func CreateDept(db *gorm.DB, data model.Dept) (err error) {
	err = db.Debug().Create(&data).Error
	return
}

// UpdateDept
//
//	@Description: 更新部门信息
//	@param db
//	@param data
//	@return err
func UpdateDept(db *gorm.DB, data model.Dept) (err error) {
	err = db.Debug().Model(&model.Dept{}).Where("id = ?", data.ID).Updates(&data).Error
	return
}

// DeleteDeptByIds
//
//	@Description: 通过ids删除部门
//	@param db
//	@param ids
//	@return err
func DeleteDeptByIds(db *gorm.DB, ids []uint) (err error) {
	err = db.Debug().Delete(&model.Dept{}, ids).Error
	return
}
