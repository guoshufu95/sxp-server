package service

import (
	"fmt"
	"gorm.io/gorm"
	"sxp-server/app/dao"
	"sxp-server/app/model"
	"sxp-server/app/service/dto"
	"time"
)

type MenuService struct {
	Service
}

// ListMenu
//
//	@Description: 返回所有菜单
//	@receiver s
//	@return err
func (s *MenuService) ListMenu() (err error, menuList []dto.ListMenusRes) {
	var menus []model.Menu
	err = dao.ListMenus(s.Db, &menus)
	if err != nil {
		s.Logger.Error("获取menus列表失败")
		return
	}
	list := make([]dto.ListMenusRes, 0)
	for _, menu := range menus {
		var l = &dto.ListMenusRes{}
		l.BuildResponse(menu)
		list = append(list, *l)
	}

	menuList = GetMenuTree0(list, 0)
	return
}

// GetRoleMenus
//
//	@Description: 返回当前用户的菜单列表
//	@receiver s
//	@return err
//	@return menus
func (s *MenuService) GetRoleMenus(name string, roleIds []int) (err error, menus []model.Menu) {
	if len(roleIds) == 0 {
		return
	}
	// admin返回所有
	if name == "admin" {
		err = dao.ListMenus(s.Db, &menus)
	} else {
		var roles []model.Role
		err = dao.RoleMenus(s.Db, roleIds, &roles)
		if err != nil {
			s.Logger.Error("获取当前用户菜单列表失败")
		}
		m := make(map[uint]model.Menu)
		for _, role := range roles { //去重
			for _, menu := range role.Menus {
				m[menu.ID] = menu
			}
		}
		for _, v := range m {
			menus = append(menus, v)
		}
	}

	menus = GetMenuTree(menus, 0)
	for i, _ := range menus {
		_, l := getMenuTreeById(s.Db, menus[i].ID)
		menus[i].Children = l
	}
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

// GetMenuTree0
//
//	@Description: 构造树状返回菜单结构
//	@param data
//	@param parentId
//	@return []dto.ListMenusRes
func GetMenuTree0(data []dto.ListMenusRes, parentId int) []dto.ListMenusRes {
	var listTree []dto.ListMenusRes
	for _, val := range data {
		if val.ParentId == parentId {
			children := GetMenuTree0(data, int(val.Id))
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

// QueryMenuByParam
//
//	@Description: 条件查询
//	@receiver s
//	@param req
func (s *MenuService) QueryMenuByParam(req dto.QueryMenusByParamReq) (err error, menuList []model.Menu) {
	db := s.buildCondition(s.Db, req.Name)
	var menus []model.Menu
	err = dao.QueryMenuByParam(db, &menus)
	if err != nil {
		s.Logger.Error("menu条件查询失败")
		return
	}
	for i, _ := range menus {
		menus[i].CreateTime = menus[i].CreatedAt.Format(time.DateTime)
	}
	if req.Name == "" {
		menuList = GetMenuTree(menus, 0)
	} else {
		var duplicateMap = make(map[uint]model.Menu)
		// 条件查询不为空时
		for _, menu := range menus {
			if _, ok := duplicateMap[menu.ID]; ok {
				continue
			}
			_, ll := getMenuTreeById(s.Db, menu.ID)
			for i, _ := range ll {
				ll[i].CreateTime = ll[i].CreatedAt.Format(time.DateTime)
			}
			getDuplicateMap(duplicateMap, ll)
			duplicateMap[menu.ID] = menu
			menu.Children = append(menu.Children, ll...)
			menuList = append(menuList, menu)
		}
	}
	return
}

// getDuplicateMap
//
//	@Description: 去重
//	@param dm
//	@param list
func getDuplicateMap(dm map[uint]model.Menu, list []model.Menu) {
	for _, m := range list {
		if len(m.Children) != 0 {
			getDuplicateMap(dm, m.Children)
		}
		dm[m.ID] = m
	}
}

// buildCondition
//
//	@Description: 构造menu条件查询sql
//	@receiver s
//	@param db
//	@param req
//	@return err
func (s *MenuService) buildCondition(db *gorm.DB, name string) *gorm.DB {
	if name != "" {
		db = db.Where(fmt.Sprintf("name like \"%s\" or name like \"%s\" or name like \"%s\" or name =\"%s\"",
			"%"+name+"%",
			"%"+name,
			name+"%",
			name))
	}
	return db
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
	var m = make(map[string]interface{})
	m["parent_id"] = menu.ParentId
	m["name"] = menu.Name
	m["path"] = menu.Path
	m["icon"] = menu.Icon
	m["order_num"] = menu.OrderNum
	m["hidden"] = menu.Hidden
	m["redirect"] = menu.Redirect
	err = dao.UpdateMenu(s.Db, menu.ID, m)
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
	for _, i := range ids {
		var m model.Menu
		_ = dao.GetMenusById(db, i, &m)
		err = dao.DeleteMenuByIds(db, m)
		if err != nil {
			s.Logger.Error("删除菜单失败")
			return
		}
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

// GetMenuById
//
//	@Description: 返回menu详情
//	@receiver s
//	@param id
//	@return err
func (s *MenuService) GetMenuById(id int) (err error, menu model.Menu) {
	err = dao.GetMenusById(s.Db, uint(id), &menu)
	if err != nil {
		s.Logger.Error("获取菜单详情失败")
		return
	}
	return
}
