package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
	logging "github.com/ipfs/go-log"
)

var (
	log = logging.Logger("web-socket")
)

type DSWebSocket struct {
	upgrader    ws.Upgrader
	connections []*ws.Conn
}

func NewDSWebSocket() *DSWebSocket {
	return &DSWebSocket{
		upgrader: ws.Upgrader{},
	}
}

// HandleFunc handles incoming ws connections.
func (dsWebSocket *DSWebSocket) HandleFunc(c *gin.Context) {

	// TODO: Check origin.
	dsWebSocket.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := dsWebSocket.upgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
	if err != nil {
		log.Error(err)
		return
	}

	// TODO: Should not just broadcast.
	dsWebSocket.connections = append(dsWebSocket.connections, conn)

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			log.Error(err)
			return
		}

		for _, c := range dsWebSocket.connections {
			go dsWebSocket.Write(c, t, msg)
		}
	}
}

func (dsWebSocket *DSWebSocket) Write(conn *ws.Conn, messageType int, data []byte) {
	err := conn.WriteMessage(messageType, data)
	if err != nil {
		log.Error(err)
	}
}
