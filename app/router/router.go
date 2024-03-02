package router

import (
	"github.com/gin-gonic/gin"
	"sxp-server/app/api/integral"
	"sxp-server/app/api/login"
	"sxp-server/app/api/product"
	"sxp-server/app/api/task"
	"sxp-server/common/middleware"
)

// InitRouter
//
//	@Description: 初始化路由
//	@param r
func InitRouter(r *gin.Engine) {
	g := r.Group("/sxp")
	g.Use(middleware.LoggerMiddleware())
	Router(g)
	//日志中间件
}

func Router(g *gin.RouterGroup) {
	buildTask(g.Group("/task"))
	buildIntegral(g.Group("/integral"))
	buildLogin(g.Group("/login"))
	buildProduct(g.Group("/product"))
}

// buildTask
//
//	@Description: 定时任务
//	@param g
func buildTask(g *gin.RouterGroup) {
	a := task.TaskApi{}
	g.POST("/start", a.StartTask)
	g.POST("/getTasks", a.GetTasks)
}

// buildIntegral
//
//	@Description: 积分
//	@param g
func buildIntegral(g *gin.RouterGroup) {
	i := integral.IntegralApi{}
	g.POST("/init", i.InitIntegral)
}

// buildLogin
//
//	@Description: 登录相关
//	@param g
func buildLogin(g *gin.RouterGroup) {
	l := login.LoginApi{}
	g.POST("/post", l.Login)
}

// buildProduct
//
//	@Description: 产品相关
//	@param g
func buildProduct(g *gin.RouterGroup) {
	p := product.ProductApi{}
	g.POST("/getProduct", p.GetProduct)
}
