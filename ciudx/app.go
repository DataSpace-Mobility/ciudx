package ciudx

import (
	redis "github.com/dataspace-mobility/rs-iudx/ciudx/redis"
	"github.com/gin-gonic/gin"
)

// App holds the app instance.
type App struct {
	// Redis connection
	RedisConnection *redis.RedisConnection
	// Router
	Router *gin.Engine
}

// NewApp creates a new instance of the app.
func NewApp() *App {
	return &App{
		RedisConnection: redis.NewRedisConnection(),
		Router:          NewRouter(),
	}
}

// Run starts the router.
func (app App) Run() error {
	return app.Router.Run(":8001")
}
