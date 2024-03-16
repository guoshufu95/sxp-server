package initial

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"sxp-server/common/cache"
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

func MakeApp() *Application {
	var app = App
	return app
}
