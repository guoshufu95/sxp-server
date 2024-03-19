package service

import "sxp-server/app/service/dto"

type RoleService struct {
	Service
}

// CreateUser
//
//	@Description: 创建user
//	@receiver s
//	@param req
//	@return err
func (s *RoleService) CreateUser(req dto.CreateUserReq) (err error) {
	// 超级管理员才能创建用户
	return
}
