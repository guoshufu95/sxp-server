package model

import "gorm.io/gorm"

// LoginReq
// @Description: 登录请求参数
type LoginReq struct {
	Username string `form:"UserName" json:"username" binding:"required"`
	Password string `form:"Password" json:"password" binding:"required"`
}

// SxpUser
// @Description: 用户数据库字段
type SxpUser struct {
	gorm.Model
	UserName string `json:"userName"` //登录名
	PassWord string `json:"passWord"` //登录密码
	NickName string `json:"nickName"` //别名
	Phone    string `json:"phone"`    //电话号码
	Sex      string `json:"sex"`      //性别
	Email    string `json:"email"`    //邮箱
	Remark   string `json:"remark"`   //备注
}
