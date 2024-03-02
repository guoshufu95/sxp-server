package dto

import "gorm.io/gorm"

type IntegralInfo struct {
	gorm.Model
	IntegralId int
	Integral   int
}

func (IntegralInfo) TableName() string {
	return "integral_info"
}
