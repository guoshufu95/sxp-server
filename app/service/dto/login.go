package dto

// LoginReq
// @Description: 登录请求参数
type LoginReq struct {
	Username string `form:"UserName" json:"username" binding:"required"`
	Password string `form:"Password" json:"password" binding:"required"`
}
