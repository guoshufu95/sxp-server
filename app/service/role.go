package service

import (
	"sxp-server/app/dao"
	"sxp-server/app/model"
	"sxp-server/app/service/dto"
)

type RoleService struct {
	Service
}

// ListRoles
//
//	@Description: 角色列表返回
//	@receiver s
//	@return err
//	@return roles
func (s *RoleService) ListRoles() (err error, roles []model.Role) {
	err = dao.ListRoles(s.Db, &roles)
	for i, _ := range roles {
		roles[i].Menus = GetMenuTree(roles[i].Menus, 0)

	}
	return
}

// CreateRole
//
//	@Description: 创建role
//	@receiver s
//	@param req
//	@return err
func (s *RoleService) CreateRole(req dto.CreateRoleReq) (err error) {
	var (
		data  model.Role
		menus []model.Menu
		detps []model.Dept
	)
	db := s.Db
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	req.BuildCreateData(&data)
	err = dao.GetMenusByIds(db, req.MenuIds, &menus)
	if err != nil {
		s.Logger.Error("通过ids查询菜单列表失败")
		return
	}
	err = dao.GetDeptsByIds(db, req.DeptIds, &detps)
	if err != nil {
		s.Logger.Error("通过ids查询部门列表失败")
		return
	}
	data.Menus = menus
	data.Depts = detps
	err = dao.CreateRole(db, data)
	if err != nil {
		s.Logger.Error("创建role失败")
		return
	}
	return
}

// UpdateRole
//
//	@Description: 更新role
//	@receiver s
//	@param req
//	@return err
func (s *RoleService) UpdateRole(req dto.UpdateRoleReq) (err error) {
	var (
		role  model.Role
		menus []model.Menu
		detps []model.Dept
	)
	db := s.Db
	db.Begin() //开启事务
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	err = dao.GetRoleById(db, req.Id, &role)
	if err != nil {
		s.Logger.Error("通过id查询role失败")
		return
	}
	req.BuildUpdateData(&role)
	err = dao.DeleteRoleDepts(db, role)
	if err != nil {
		s.Logger.Error("删除角色部门失败")
		return
	}
	err = dao.DeleteRoleMenus(db, role)
	if err != nil {
		s.Logger.Error("删除角色菜单失败")
		return
	}
	err = dao.GetMenusByIds(db, req.MenuIds, &menus)
	if err != nil {
		s.Logger.Error("通过ids查询菜单列表失败")
		return
	}
	err = dao.GetDeptsByIds(db, req.DeptIds, &detps)
	if err != nil {
		s.Logger.Error("通过ids查询部门列表失败")
		return
	}
	role.Menus = menus
	role.Depts = detps
	err = dao.UpdateRole(db, role)
	if err != nil {
		s.Logger.Error("更新role失败")
	}
	return
}
