package dao

import (
	"gorm.io/gorm"
	"sxp-server/app/dao/dto"
)

// InertIntegralInfo
//
//	@Description: 积分初始化字段入库
//	@param db
//	@param id
//	@param integral
//	@return err
func InertIntegralInfo(db *gorm.DB, id, integral int) (err error) {
	var param = dto.IntegralInfo{
		IntegralId: id,
		Integral:   integral,
	}
	return db.Debug().Create(&param).Error
}
