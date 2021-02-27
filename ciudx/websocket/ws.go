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
	connections map[string](map[string]*ws.Conn)
}

func NewDSWebSocket() *DSWebSocket {
	return &DSWebSocket{
		upgrader:    ws.Upgrader{},
		connections: make(map[string](map[string]*ws.Conn)),
	}
}

// HandleFunc handles incoming ws connections.
func (dsWebSocket *DSWebSocket) HandleFunc(c *gin.Context) {

	resource, ok := c.Params.Get("resource")
	if !ok {
		log.Info("Resource not provided")
		return
	}
	perm, ok := c.Params.Get("perm")
	if !ok {
		log.Info("perm not provided")
		return
	}

	// TODO: Check origin.
	dsWebSocket.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := dsWebSocket.upgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
	if err != nil {
		log.Error(err)
		return
	}

	if perm == "write" {
		if conMap, ok := dsWebSocket.connections[resource]; ok {
			for {
				t, msg, err := conn.ReadMessage()
				if err != nil {
					log.Error(err)
					return
				}

				for _, con := range conMap {
					go dsWebSocket.Write(con, t, msg)
				}
			}
		}

	} else if perm == "read" {
		if _, ok := dsWebSocket.connections[resource]; !ok {
			dsWebSocket.connections[resource] = make(map[string]*ws.Conn)
			dsWebSocket.connections[resource][conn.RemoteAddr().String()] = conn
		}
	}

	// TODO: Should not just broadcast.
	// dsWebSocket.connections = append(dsWebSocket.connections, conn)
	// log.Info(dsWebSocket.connections)

}

// Send sends the data to all conns.
// func (dsWebSocket *DSWebSocket) Send(data []byte) error {
// 	for _, c := range dsWebSocket.connections {
// 		return dsWebSocket.Write(c, ws.TextMessage, data)
// 	}
// 	return nil
// }

func (dsWebSocket *DSWebSocket) Write(conn *ws.Conn, messageType int, data []byte) error {
	err := conn.WriteMessage(messageType, data)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
