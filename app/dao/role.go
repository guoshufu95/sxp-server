package dao

import (
	"gorm.io/gorm"
	"sxp-server/app/model"
)

func GetRoleById(db *gorm.DB, id int, role *model.Role) (err error) {
	err = db.Table("role").Debug().Where("id = ?", id).Find(&role).Error
	return
}
