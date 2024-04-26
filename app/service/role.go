package service

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
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
	req.BuildUpdateData(&role)
	role.Menus = menus
	role.Depts = detps
	err = dao.UpdateRoleMenus(db, role)
	if err != nil {
		s.Logger.Error("更新关联菜单失败")
		return
	}
	err = dao.UpdateRoleDepts(db, role)
	if err != nil {
		s.Logger.Error("更新关联部门失败")
		return
	}
	err = dao.UpdateRole(db, role)
	if err != nil {
		s.Logger.Error("更新role失败")
	}
	return
}

// DeleteRole
//
//	@Description: 删除用户
//	@receiver s
//	@param id
//	@return err
func (s *RoleService) DeleteRole(id int) (err error) {
	db := s.Db
	db.Begin()
	defer func() {
		if err != nil {
			db.Callback()
		} else {
			db.Commit()
		}
	}()
	var role model.Role
	err = dao.GetRoleById(db, id, &role)
	if err != nil {
		s.Logger.Error("通过id查询role失败")
		return
	}
	err = dao.DeleteRoleById(db, role)
	if err != nil {
		s.Logger.Error("删除role失败")
		return
	}
	return
}

// UpdateStatus
//
//	@Description: 更新启用状态
//	@receiver s
//	@param req
//	@return err
func (s *RoleService) UpdateStatus(req dto.UpdateRoleStatusReq) (err error) {
	err = dao.UpdateRoleStatus(s.Db, req.Id, req.Status)
	if err != nil {
		s.Logger.Error("更新启用状态失败")
		return
	}

	return
}

// RoleByParams
//
//	@Description: 角色列表条件查询
//	@receiver s
//	@param req
//	@return err
func (s *RoleService) RoleByParams(req dto.QueryRoleByParams) (err error, roles []model.Role) {
	db := s.buildCondition(s.Db, req)
	err = dao.GetRoleByParams(db, &roles)
	if err != nil {
		s.Logger.Error("通过参数查询roles失败")
		return
	}
	return
}

// GetRoleById
//
//	@Description: 通过id查询role详情返回
//	@receiver s
//	@param id
//	@return err
//	@return role
func (s *RoleService) GetRoleById(id int) (err error, res dto.GetRoleByIdRes) {
	var role model.Role
	err = dao.GetRoleById(s.Db, id, &role)
	if err != nil {
		s.Logger.Error("id查询role返回失败")
		return
	}
	res.Id = role.ID
	res.Name = role.Name
	res.Status = role.Status
	res.RoleKey = role.RoleKey
	res.Label = role.Label
	res.RoleSort = strconv.Itoa(role.RoleSort)
	for _, menu := range role.Menus {
		res.MenuIds = append(res.MenuIds, menu.ID)
	}
	for _, dept := range role.Depts {
		res.DeptIds = append(res.DeptIds, dept.ID)
	}
	return
}

// buildCondition
//
//	@Description: 构造条件查询语句
//	@receiver s
//	@param db
//	@param req
//	@return *gorm.DB
func (s *RoleService) buildCondition(db *gorm.DB, req dto.QueryRoleByParams) *gorm.DB {
	if req.Name != "" {
		db = db.Where(fmt.Sprintf("name like \"%s\" or name like \"%s\" or name like \"%s\" or name =\"%s\"",
			"%"+req.Name+"%",
			"%"+req.Name,
			req.Name+"%",
			req.Name))
	}
	if req.RoleKey != "" {
		db = db.Where(fmt.Sprintf("role_key like \"%s\" or role_key like \"%s\" or role_key like \"%s\" or role_key =\"%s\"",
			"%"+req.RoleKey+"%",
			"%"+req.RoleKey,
			req.RoleKey+"%",
			req.RoleKey))
	}
	if req.Status != "" {
		if req.Status == "启用" {
			db = db.Where("status = ?", 1)
		} else {
			db = db.Where("status = ?", 0)
		}
	}
	return db
}
