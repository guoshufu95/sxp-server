package service

import (
	"gorm.io/gorm"
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
	err = dao.CreateMenu(s.Db, menu)
	if err != nil {
		s.Logger.Error("创建菜单入库口失败")
		return
	}
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
	err = dao.UpdateMenu(s.Db, menu)
	if err != nil {
		s.Logger.Error("更新菜单失败")
		return
	}
	return
}

// DeleteMenu
//
//	@Description:
//	@receiver s
//	@param id
//	@return err
func (s *MenuService) DeleteMenu(id int) (err error) {
	var menu model.Menu
	db := s.Db
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	err = dao.GetMenusById(db, uint(id), &menu)
	if err != nil {
		s.Logger.Error("通过id查询menu失败")
		return
	}
	err, list := getMenuTreeById(db, menu.ID)
	if err != nil {
		s.Logger.Error("获取菜单树错误")
		return
	}
	menu.Children = append(menu.Children, list...)
	ids := make([]uint, 0)
	getMenuTreeIds(menu.Children, &ids)
	ids = append(ids, menu.ID)
	err = dao.DeleteMenuByIds(db, ids)
	if err != nil {
		s.Logger.Error("删除菜单失败")
		return
	}
	return
}

// getMenuTreeById
//
//	@Description: 获取该id下的菜单树
//	@param db
//	@param id
//	@return err
//	@return list
func getMenuTreeById(db *gorm.DB, id uint) (err error, list []model.Menu) {
	var menus []model.Menu
	err = dao.ListMenus(db, &menus)
	if err != nil {
		return
	}
	list = GetMenuTree(menus, int(id))
	return
}

// getMenuTreeIds
//
//	@Description: 返回children ids
//	@param menus
//	@param ids
func getMenuTreeIds(menus []model.Menu, ids *[]uint) {
	for _, val := range menus {
		*ids = append(*ids, val.ID)
		if len(val.Children) != 0 {
			getMenuTreeIds(val.Children, ids)
		}
	}
}
