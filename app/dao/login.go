package dao

import (
	"gorm.io/gorm"
	"sxp-server/app/model"
)

// GetUser
//
//	@Description: 根据用户名查询
//	@param db
//	@param name
//	@return err
func GetUser(db *gorm.DB, name string) (err error, user model.SxpUser) {
	err = db.Table("sxp_user").Where("name = ?", name).Find(&user).Error
	return
}
