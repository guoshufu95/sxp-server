package service

import (
	"errors"
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
func (s *DeptService) GetDept() (err error, dept model.Dept) {
	var depts []model.Dept
	err = dao.GetAllDepts(s.Db, &depts)
	if err != nil {
		s.Logger.Error("查询所有部门列表失败")
		return
	}
	dept = getDeptTree(depts, 0)[0]
	return
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

// UpdateDept
//
//	@Description: 更新部门信息
//	@receiver s
//	@param req
//	@return err
func (s *DeptService) UpdateDept(req dto.UpdateDeptReq) (err error) {
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
	if dm.ID != uint(req.Id) {
		err = errors.New("部门名重复，请重新设置")
		return
	}
	req.BuildUpdateData(&dept)
	err = dao.UpdateDept(db, dept)
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
	ids := make([]uint, 0)
	getDeptTreeIds(dept.Children, &ids)
	ids = append(ids, dept.ID)
	err = dao.DeleteDeptByIds(db, ids)
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
func getDeptTreeIds(depts []model.Dept, ids *[]uint) {
	for _, val := range depts {
		*ids = append(*ids, val.ID)
		if len(val.Children) != 0 {
			getDeptTreeIds(val.Children, ids)
		}
	}
}
