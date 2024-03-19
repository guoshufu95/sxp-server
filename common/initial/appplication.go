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

type Application struct {
	ProjectName string         `json:"projectName"`
	Engine      *gin.Engine    `json:"engine"`
	Db          *gorm.DB       `json:"globalDb"`
	Cache       *redis.Client  `json:"cache"`
	Logger      *logger.ZapLog `json:"logger"`
	mux         sync.Mutex
	Casbin      *casbin.SyncedEnforcer `json:"casbins"`
}

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

func (a *Application) GetAppDb() *gorm.DB {
	a.mux.Lock()
	defer a.mux.Unlock()
	return a.Db
}

func (a *Application) GetCache() *redis.Client {
	a.mux.Lock()
	defer a.mux.Unlock()
	return a.Cache
}

func (a *Application) GetCasbin() *casbin.SyncedEnforcer {
	a.mux.Lock()
	defer a.mux.Unlock()
	return a.Casbin
}
