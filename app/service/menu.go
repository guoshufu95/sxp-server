package service

import (
	"sxp-server/app/model"
	"sxp-server/app/service/dto"
)

type MenuService struct {
	Service
}

// CreateMenu
//
//	@Description: 创建菜单
//	@receiver s
//	@param req
//	@return err
func (s *MenuService) CreateMenu(req dto.CreateMenuReq) (err error) {
	var menu model.Menu
	req.BuildCreateData(&menu)
	err = s.Db.Create(&menu).Error
	return
}

// UpdateMenu
//
//	@Description: 更新
//	@receiver s
//	@param req
//	@return err
func (s *MenuService) UpdateMenu(req dto.UpdateMenuReq) (err error) {
	var menu model.Menu
	req.BuildUpdateData(&menu)
	err = s.Db.Model(&menu).Updates(&menu).Error
	return
}
