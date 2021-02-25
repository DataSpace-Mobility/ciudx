package ciudx

import (
	"os"

	redis "github.com/dataspace-mobility/rs-iudx/ciudx/redis"
	"github.com/dataspace-mobility/rs-iudx/ciudx/utils"
	"github.com/dataspace-mobility/rs-iudx/ciudx/websocket"
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

	ws := websocket.NewDSWebSocket()
	router := NewRouter(app, ws)
	app.Router = router

	return app
}

// Run starts the router.
func (app App) Run() error {
	port := utils.Getenv("LISTEN_PORT", "8001")
	log.Info("Starting RS-IUDX app on port ", port)
	return app.Router.Run(":" + port)
}

func runDebug() {
	os.Setenv("REDIS_HOST", "localhost")
}
