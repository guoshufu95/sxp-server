package initial

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"sxp-server/common/cache"
	"sxp-server/common/db"
	"sxp-server/common/logger"
	"sxp-server/config"
)

var App *Application

type Application struct {
	ProjectName string         `json:"projectName"`
	Engine      *gin.Engine    `json:"engine"`
	Db          *gorm.DB       `json:"globalDb"`
	Cache       *redis.Client  `json:"cache"`
	Logger      *logger.ZapLog `json:"logger"`
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

func MakeApp() *Application {
	var app = App
	return app
}
