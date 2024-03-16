package router

import (
	"github.com/gin-gonic/gin"
	"sxp-server/app/api"
	"sxp-server/common/initial"
	"sxp-server/common/middleware"
)

// InitRouter
//
//	@Description: 初始化路由
//	@param r
func InitRouter(r *gin.Engine) {
	g := r.Group("/sxp")
	g.Use(middleware.LoggerMiddleware()).
		Use(middleware.WithGormDb(initial.App.GetAppDb())).
		Use(middleware.WithRedisDb(initial.App.GetCache()))
	Router(g)
	//日志中间件
}

func Router(g *gin.RouterGroup) {
	buildTask(g.Group("/task"))
	buildIntegral(g.Group("/integral"))
	buildLogin(g.Group("/login"))
	buildProduct(g.Group("/product"))
	buildMenu(g.Group("/menu"))
}

// buildTask
//
//	@Description: 定时任务路由
//	@param g
func buildTask(g *gin.RouterGroup) {
	a := api.TaskApi{}
	g.Use(middleware.JWTAuthMiddleware())
	g.POST("/start", a.StartTask)
	g.POST("/getTasks", a.GetTasks)
}

// buildIntegral
//
//	@Description: 积分功能路由
//	@param g
func buildIntegral(g *gin.RouterGroup) {
	i := api.IntegralApi{}
	g.POST("/init", i.InitIntegral)
}

// buildLogin
//
//	@Description: 登录路由
//	@param g
func buildLogin(g *gin.RouterGroup) {
	l := api.LoginApi{}
	g.POST("/", l.Login)
}

// buildProduct
//
//	@Description: 产品路由
//	@param g
func buildProduct(g *gin.RouterGroup) {
	g.Use(middleware.JWTAuthMiddleware())
	p := api.ProductApi{}
	g.POST("/getProduct", p.GetProduct)
	g.POST("/updateProduct", p.UpdateProduct)
	g.POST("/getByStatus", p.GetByStatus)
}

// buildMenu
//
//	@Description: 菜单路由
//	@param g
func buildMenu(g *gin.RouterGroup) {
	g.Use(middleware.JWTAuthMiddleware())
	m := api.MenuApi{}
	g.POST("/get", m.GetMenu)
	g.POST("/create", m.CreateMenu)
	g.POST("/update", m.UpdateMenu)
	g.POST("/delete", m.DeleteMenu)
}
