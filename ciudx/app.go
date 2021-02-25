package ciudx

import (
	"os"

	redis "github.com/dataspace-mobility/rs-iudx/ciudx/redis"
	"github.com/gin-gonic/gin"
	logging "github.com/ipfs/go-log"
)

var (
	log = logging.Logger("app")
)

// App holds the app instance.
type App struct {
	// Redis connection
	RedisConnection *redis.Connection
	// Router
	Router *gin.Engine
}

// NewApp creates a new instance of the app.
func NewApp() *App {
	runDebug()

	app := &App{
		RedisConnection: redis.NewRedisConnection(),
	}

	router := NewRouter(app)
	app.Router = router

	return app
}

// Run starts the router.
func (app App) Run() error {
	log.Info("Starting RS-IUDX app on port ", "8001")
	return app.Router.Run(":8001")
}

func runDebug() {
	os.Setenv("REDIS_HOST", "localhost")
}
