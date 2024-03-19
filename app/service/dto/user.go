package dto

import "sxp-server/app/model"

// CreateUserReq
// @Description:  创建角色入参
type CreateUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	NickName string `json:"nick_name"`
	Sex      string `json:"sex"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Remark   string `json:"remark"`
	Status   string `json:"status"`
	RoleId   int    `json:"roleId"`
	IsSuper  int    `json:"is_super"`
}

// BuildData
//
//	@Description: 构造入库字段
//	@receiver c
//	@param user
func (c CreateUserReq) BuildData(user *model.User) {
	user.Username = c.Username
	user.Password = c.Password
	user.NickName = c.NickName
	user.Sex = c.Sex
	user.Email = c.Email
	user.Phone = c.Phone
	user.Remark = c.Remark
	user.Status = c.Status
	user.RoleId = c.RoleId
	user.IsSuper = c.IsSuper
	return
}

// GetById
// @Description: id查询入参
type GetUserByIdRequest struct {
	Id int `json:"id"`
}
