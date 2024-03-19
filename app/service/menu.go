package service

import (
	"sxp-server/app/dao"
	"sxp-server/app/model"
	"sxp-server/app/service/dto"
)

type MenuService struct {
	Service
}

// ListMenu
//
//	@Description: 返回所有菜单
//	@receiver s
//	@return err
func (s *MenuService) ListMenu() (err error, menus []model.Menu) {
	err = dao.ListMenus(s.Db, &menus)
	if err != nil {
		s.Logger.Error("获取menus列表失败")
		return
	}
	menus = GetMenuTree(menus, 0)
	return
}

// GetRoleMenus
//
//	@Description: 返回当前用户的菜单列表
//	@receiver s
//	@return err
//	@return menus
func (s *MenuService) GetRoleMenus(roleKey string, roleId int) (err error, menus []model.Menu) {
	// admin返回所有
	if roleKey == "admin" {
		err = dao.ListMenus(s.Db, &menus)
	} else {
		var role model.Role
		err = dao.RoleMenus(s.Db, roleId, &role)
		if err != nil {
			s.Logger.Error("获取当前用户菜单列表失败")
		}
		menus = role.Menus
	}
	menus = GetMenuTree(menus, 0)
	return
}

// GetMenuTree
//
//	@Description: 递归生成菜单结构
//	@param data
//	@param parentId
//	@return []model.Menu
func GetMenuTree(data []model.Menu, parentId int) []model.Menu {
	var listTree []model.Menu
	for _, val := range data {
		if val.ParentId == parentId {
			children := GetMenuTree(data, int(val.ID))
			if len(children) > 0 {
				val.Children = children
			}
			listTree = append(listTree, val)
		}
	}
	return listTree
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
