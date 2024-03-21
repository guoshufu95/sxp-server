package initial

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"sxp-server/common/cache"
	cs "sxp-server/common/casbin"
	"sxp-server/common/db"
	"sxp-server/common/logger"
	"sxp-server/config"
	"sync"
)

var App *Application

// Application
// @Description: 全局application
type Application struct {
	ProjectName string         `json:"projectName"`
	Engine      *gin.Engine    `json:"engine"`
	Db          *gorm.DB       `json:"globalDb"`
	Cache       *redis.Client  `json:"cache"`
	Logger      *logger.ZapLog `json:"logger"`
	mux         sync.Mutex
	Casbin      *casbin.SyncedEnforcer `json:"casbins"`
}

// init
//
//	@Description: 初始化
func init() {
	config.ReadConfig("./config/sxp.yml")
	App = &Application{
		Logger:      logger.GetLogger(),
		ProjectName: "sxp-server",
		Engine:      gin.Default(),
		Db:          db.IniDb(),
		Cache:       cache.IniCache(),
	}
	e := cs.InitCabin(App.Db)
	App.Casbin = e
}

// GetAppDb
//
//	@Description: 返回全局的db
//	@receiver a
//	@return *gorm.DB
func (a *Application) GetAppDb() *gorm.DB {
	a.mux.Lock()
	defer a.mux.Unlock()
	return a.Db
}

// GetCache
//
//	@Description: 返回一个全局的rdb
//	@receiver a
//	@return *redis.Client
func (a *Application) GetCache() *redis.Client {
	a.mux.Lock()
	defer a.mux.Unlock()
	return a.Cache
}

// GetCasbin
//
//	@Description: 返回一个全局的casbin SyncedEnforcer
//	@receiver a
//	@return *casbin.SyncedEnforcer
func (a *Application) GetCasbin() *casbin.SyncedEnforcer {
	a.mux.Lock()
	defer a.mux.Unlock()
	return a.Casbin
}
