package model

import "gorm.io/gorm"

type IntegralInfo struct {
	gorm.Model
	IntegralsId int `json:"integralId" gorm:"type:int(10);"`
	Integrals   int `json:"Integrals" gorm:"type:int(100)"`
}

func (IntegralInfo) TableName() string {
	return "integral_info"
}
