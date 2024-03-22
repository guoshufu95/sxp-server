package dto

import "sxp-server/app/model"

type CommonUserReq struct {
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

func (c CommonUserReq) buildData(user *model.User) {
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
}

// CreateUserReq
// @Description:  创建角色入参
type CreateUserReq struct {
	CommonUserReq
}

// BuildCreateData
//
//	@Description: 构造入库字段
//	@receiver c
//	@param user
func (c CreateUserReq) BuildCreateData(user *model.User) {
	c.CommonUserReq.buildData(user)
}

// UpdateUserReq
// @Description: 更新数据入参
type UpdateUserReq struct {
	Id int `json:"id"`
	CommonUserReq
}

// BuildUpdateData
//
//	@Description: 构造更新参数
//	@receiver c
//	@param user
func (c UpdateUserReq) BuildUpdateData(user *model.User) {
	c.CommonUserReq.buildData(user)
	user.ID = uint(c.Id)
}

// GetUserByIdRequest
// @Description: id查询入参
type GetUserByIdRequest struct {
	Id int `json:"id"`
}

// DeleteUserReq
// @Description:删除user入参
type DeleteUserReq struct {
	Id int `json:"id"`
}
