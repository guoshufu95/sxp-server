package dao

import (
	"gorm.io/gorm"
	"sxp-server/app/model"
)

// GetDeptsByIds
//
//	@Description: 返回部门列表
//	@param db
//	@param ids
//	@param depts
//	@return err
func GetDeptsByIds(db *gorm.DB, ids []int, depts *[]model.Dept) (err error) {
	err = db.Debug().Find(&depts, ids).Error
	return
}
