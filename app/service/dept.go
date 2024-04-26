package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sxp-server/app/dao"
	"sxp-server/app/model"
	"sxp-server/app/service/dto"
)

type DeptService struct {
	Service
}

// GetDept
//
//	@Description: 返回部门信息
//	@receiver s
//	@return err
//	@return dept
func (s *DeptService) GetDept() (err error, res []dto.DeptsTree) {
	var depts []model.Dept
	err = dao.GetAllDepts(s.Db, &depts)
	if err != nil {
		s.Logger.Error("查询所有部门列表失败")
		return
	}
	var tree []dto.DeptsTree
	dto.BuildDeptsTreeRes(depts, &tree)
	//dept = getDeptTree(depts, 0)[0]
	res = GetTree(tree, 0)
	return
}

// GetDeptByParams
//
//	@Description: 条件查询
//	@receiver s
//	@param name
//	@return err
func (s *DeptService) GetDeptByParams(req dto.DeptByParamsReq) (err error, res []dto.DeptsTree) {
	var (
		depts = make([]model.Dept, 0)
		dps   = make([]model.Dept, 0)
	)

	if req.Name == "" && req.Status == "" {
		err = dao.GetAllDepts(s.Db, &depts)
		if err != nil {
			s.Logger.Error("查询失败")
			return
		}
	} else {
		db := s.buildDeptQuery(s.Db, req)
		err = dao.GetDeptsByParams(db, &dps)
		for _, dept := range dps {
			if dept.ParentId != 0 {
				var m = make(map[uint]model.Dept)
				s.buildParentDept(dept, m)
				for _, v := range m {
					depts = append(depts, v)
				}
			}
		}
	}
	dm := make(map[uint]model.Dept)
	for _, v := range depts {
		dm[v.ID] = v
	}
	// 去重
	var list = make([]model.Dept, 0)
	for _, v := range dm {
		list = append(list, v)
	}
	var tree []dto.DeptsTree
	dto.BuildDeptsTreeRes(list, &tree)
	res = GetTree(tree, 0)
	return
}

// buildDeptQuery
//
//	@Description: 构造dept条件查询参数
//	@receiver s
//	@param req
//	@return err
func (s *DeptService) buildDeptQuery(db *gorm.DB, req dto.DeptByParamsReq) *gorm.DB {
	if req.Name != "" {
		db = db.Where(fmt.Sprintf("name like \"%s\" or name like \"%s\" or name like \"%s\" or name =\"%s\"",
			"%"+req.Name+"%",
			"%"+req.Name,
			req.Name+"%",
			req.Name))
	}
	if req.Status != "" {
		if req.Status == "正常" {
			db = db.Where("status = ?", 1)
		} else {
			db = db.Where("status = ?", 0)
		}
	}
	return db
}

// buildParentDept
//
//	@Description: 组装父节点部门信息返回
//	@receiver s
//	@param dept
//	@return res
func (s *DeptService) buildParentDept(dept model.Dept, m map[uint]model.Dept) {
	var (
		dp model.Dept
	)
	if dept.ParentId != 0 {
		err := dao.GetDeptById(s.Db, dept.ParentId, &dp)
		if err != nil {
			s.Logger.Error("查询部门信息失败!")
			return
		}
		m[dept.ID] = dept
		s.buildParentDept(dp, m)
		return
	} else {
		err := dao.GetDeptById(s.Db, 1, &dp)
		if err != nil {
			s.Logger.Error("查询失败")
			return
		}
		m[dp.ID] = dp
		return
	}
}

// getDeptTree
//
//	@Description: 构建dept tree列表
//	@param data
//	@param parentId
//	@return []model.Dept
func getDeptTree(data []model.Dept, parentId uint) []model.Dept {
	var listTree []model.Dept
	for _, val := range data {
		if val.ParentId == parentId {
			children := getDeptTree(data, val.ID)
			if len(children) > 0 {
				val.Children = children
			}
			listTree = append(listTree, val)
		}
	}
	return listTree
}

func GetTree(data []dto.DeptsTree, parentId uint) []dto.DeptsTree {
	var listTree []dto.DeptsTree
	for _, val := range data {
		if val.ParentId == parentId {
			children := GetTree(data, val.Id)
			if len(children) > 0 {
				val.Children = children
			}
			listTree = append(listTree, val)
		}
	}
	return listTree
}

// CreateDept
//
//	@Description: 创建部门
//	@receiver s
//	@param req
//	@return err
func (s *DeptService) CreateDept(req dto.CreateDeptReq) (err error) {
	var dept model.Dept
	db := s.Db
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	err, dm := dao.GetDeptByName(db, req.Name)
	if err != nil {
		s.Logger.Error("通过名称查询部门失败")
		return
	}
	if dm.ID != 0 {
		err = errors.New("部门名重复，请重新设置")
		return
	}
	req.BuildCreateData(&dept)
	err = dao.CreateDept(db, dept)
	if err != nil {
		s.Logger.Error("创建部门入库失败")
		return
	}
	return
}

// GetById
//
//	@Description: 通过id查询dept
//	@receiver s
//	@param id
//	@return err
//	@return dept
func (s *DeptService) GetById(id uint) (err error, dept model.Dept) {
	err = dao.GetDeptById(s.Db, id, &dept)
	if err != nil {
		s.Logger.Error("通过id返回失败")
	}
	return
}

// UpdateDept
//
//	@Description: 更新部门信息
//	@receiver s
//	@param req
//	@return err
func (s *DeptService) UpdateDept(req dto.UpdateDeptReq) (err error) {
	db := s.Db
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	err, dm := dao.GetDeptByName(db, req.Name)
	if err != nil {
		s.Logger.Error("通过名称查询部门失败")
		return
	}
	if dm.ID != 0 && dm.ID != uint(req.Id) {
		err = errors.New("部门名重复，请重新设置")
		return
	}
	m := make(map[string]interface{})
	m["id"] = req.Id
	m["parent_id"] = req.ParentId
	m["name"] = req.Name
	m["leader"] = req.Leader
	m["phone"] = req.Phone
	m["email"] = req.Email
	m["status"] = req.Status
	err = dao.UpdateDept(db, m)
	if err != nil {
		s.Logger.Error("更新部门信息失败")
		return
	}
	return
}

// DeleteDept
//
//	@Description: 删除部门信息
//	@receiver s
//	@param id
//	@return err
func (s *DeptService) DeleteDept(id int) (err error) {
	var dept model.Dept
	db := s.Db
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	err = dao.GetDeptById(db, uint(id), &dept)
	if err != nil {
		s.Logger.Error("通过id查询部门信息失败")
		return
	}
	err, list := getDeptTreeById(db, id)
	if err != nil {
		return
	}
	dept.Children = append(dept.Children, list...)
	ids := make([]int, 0)
	getDeptTreeIds(dept.Children, &ids)
	ids = append(ids, int(dept.ID))
	var depts []model.Dept
	err = dao.GetDeptsByIds(db, ids, &depts)
	if err != nil {
		s.Logger.Error("通过ids查询失败")
		return
	}
	err = dao.DeleteDeptByIds(db, depts)
	if err != nil {
		s.Logger.Error("通过ids删除部门失败")
	}
	return
}

// getDeptTreeById
//
//	@Description: 循环获取子菜单
//	@param db
//	@param id
//	@return err
//	@return list
func getDeptTreeById(db *gorm.DB, id int) (err error, list []model.Dept) {
	var depts []model.Dept
	err = dao.GetAllDepts(db, &depts)
	if err != nil {
		return
	}
	list = getDeptTree(depts, uint(id))
	return
}

// getDeptTreeIds
//
//	@Description: 返回id列表
//	@param depts
//	@return []uint
func getDeptTreeIds(depts []model.Dept, ids *[]int) {
	for _, val := range depts {
		*ids = append(*ids, int(val.ID))
		if len(val.Children) != 0 {
			getDeptTreeIds(val.Children, ids)
		}
	}
}
