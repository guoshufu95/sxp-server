package dto

import (
	"strconv"
	"sxp-server/app/model"
	"time"
)

// CommonUserReq
// @Description: 新增用户的req
type CommonUserReq struct {
	Username string `json:"username" binding:"min=3,max=255"`
	Password string `json:"password" binding:"min=3,max=20"`
	NickName string `json:"nick_name" binding:"required"`
	Sex      string `json:"sex" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Remark   string `json:"remark"`
	Status   string `json:"status"`
	DeptIds  []int  `json:"depts"`
	IsSuper  int    `json:"is_super"`
}

// UpdateUserReq
// @Description: 更新用户的req
type UpdateUserReq struct {
	Id       int    `json:"id"`
	Username string `json:"username" binding:"min=3,max=255"`
	NickName string `json:"nick_name" binding:"required"`
	Sex      string `json:"sex" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Remark   string `json:"remark"`
	Status   string `json:"status"`
	DeptIds  []int  `json:"depts"`
	IsSuper  int    `json:"is_super"`
}

// UpdateStatusReq
// @Description: 修改用户上下线状态
type UpdateStatusReq struct {
	Id     int `json:"id"`
	Status int `json:"status"`
}

// BuildUpdateData
//
//	@Description: 构建更新数据
//	@receiver c
//	@param user
func (c UpdateUserReq) BuildUpdateData(user *model.User) {
	user.ID = uint(c.Id)
	user.Username = c.Username
	user.NickName = c.NickName
	user.Sex = c.Sex
	user.Email = c.Email
	user.Phone = c.Phone
	user.Remark = c.Remark
	if c.Status != "" {
		status, _ := strconv.Atoi(c.Status)
		user.Status = status
	}
	user.IsSuper = c.IsSuper
}

func (c CommonUserReq) buildData(user *model.User) {
	user.Username = c.Username
	user.Password = c.Password
	user.NickName = c.NickName
	user.Sex = c.Sex
	user.Email = c.Email
	user.Phone = c.Phone
	user.Remark = c.Remark
	if c.Status != "" {
		status, _ := strconv.Atoi(c.Status)
		user.Status = status
	}
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

// QueryByParamsReq
// @Description: 条件查询参数
type QueryByParamsReq struct {
	UserName string `json:"username"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
}

// QueryRes
// @Description: 查询返回参数
type QueryRes struct {
	Id            uint         `json:"id"`
	Username      string       `json:"username"`
	Password      string       `json:"password"`
	NickName      string       `json:"nick_name"`
	Sex           string       `json:"sex"`
	Email         string       `json:"email"`
	Phone         string       `json:"phone"`
	LoginType     int          `json:"login_type"`
	LastLoginTime string       `json:"last_login_time"`
	Remark        string       `json:"remark"`
	Status        string       `json:"status"`
	IsSuper       int          `json:"is_super"`
	Depts         []model.Dept `json:"depts"`
}

type QueryRes0 struct {
	Id            uint   `json:"id"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	NickName      string `json:"nick_name"`
	Sex           string `json:"sex"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	LoginType     int    `json:"login_type"`
	LastLoginTime string `json:"last_login_time"`
	Remark        string `json:"remark"`
	Status        string `json:"status"`
	IsSuper       int    `json:"is_super"`
	Depts         []uint `json:"depts"`
}

// BuildQueryRes
//
//	@Description: 构建返回res
//	@param users
//	@param res
func BuildQueryRes(users *[]model.User, res *[]QueryRes) {
	for _, user := range *users {
		var status string
		if user.Status == 1 {
			status = "在线"
		} else {
			status = "下线"
		}
		var t string
		if user.LastLoginTime != nil {
			t = user.LastLoginTime.Format(time.DateTime)
		}
		var v = QueryRes{
			Id:            user.ID,
			Username:      user.Username,
			Password:      user.Password,
			NickName:      user.NickName,
			Sex:           user.Sex,
			Email:         user.Email,
			Phone:         user.Phone,
			LastLoginTime: t,
			Remark:        user.Remark,
			Status:        status,
			IsSuper:       user.IsSuper,
			Depts:         user.Depts,
		}
		*res = append(*res, v)
	}
	return
}
