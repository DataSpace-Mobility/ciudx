package video

import (
	"context"
	"io"
	"net/http"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
	logging "github.com/ipfs/go-log"
)

var (
	log = logging.Logger("video-server")
)

type VideoServer struct {
	upgrader ws.Upgrader
	cmdMap   map[string]chan []byte
}

func NewVideoServer() *VideoServer {
	return &VideoServer{
		upgrader: ws.Upgrader{
			ReadBufferSize:  10485760,
			WriteBufferSize: 10485760,
		},
		cmdMap: make(map[string](chan []byte), 15),
	}
}

// HandleFunc handles incoming ws connections.
func (videoServer *VideoServer) HandleFunc(c *gin.Context) {
	/* buf := make([]byte, 10485760)
	for {
		c.Request.Body.Read(buf)

		//log.Info(num)
	} */
	type Request struct {
		Resource string
	}

	var request Request

	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Request body is not JSON",
		})
		return
	}

	// if cmd, ok := videoServer.cmdMap[request.Resource]; ok {

	// } else {
	// 	go videoServer.runFFMPEG(c)
	// }

	// videoServer.runFFMPEG(c)
}

// WSHandleFunc handles incoming ws connections for video output.
func (videoServer *VideoServer) WSHandleFunc(c *gin.Context) {
	var resource string
	if res, ok := c.Params.Get("resource"); ok {
		resource = res
	} else {
		log.Error("no resource provided")
		return
	}

	var streamOutput chan []byte
	if chanVideoBuf, ok := videoServer.cmdMap[resource]; ok {
		streamOutput = chanVideoBuf
	} else {
		streamOutput = make(chan []byte, 10)
		videoServer.cmdMap[resource] = streamOutput

		go videoServer.runFFMPEG(resource, streamOutput)
	}

	// TODO: Check origin.
	videoServer.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := videoServer.upgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
	if err != nil {
		log.Error(err)
		return
	}

	go videoServer.copyFFMpegToWebSocket(streamOutput, conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func (videoServer *VideoServer) removeResource(resource string) {
	delete(videoServer.cmdMap, resource)
}

func (videoServer *VideoServer) copyFFMpegToWebSocket(inStream chan []byte, outConn *ws.Conn) {
	for {
		data, ok := <-inStream
		if !ok {
			log.Info("Channel Closed")
			return
		}
		err := outConn.WriteMessage(ws.BinaryMessage, data)
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func (videoServer *VideoServer) runFFMPEG(resource string, outChan chan []byte) {
	rtsp := "rtsp://localhost:8554/atcs"
	params := []string{
		"-rtsp_transport",
		"tcp",
		"-re",
		"-i",
		rtsp,
		"-q",
		"5",
		"-f",
		"mpegts",
		"-fflags",
		"nobuffer",
		"-c:v",
		"mpeg1video",
		"-an",
		"-s",
		"426x240",
		"-",
	}

	ctx, _ := context.WithTimeout(context.Background(), 200*time.Second)
	cmd := exec.CommandContext(ctx, "ffmpeg", params...)
	cmd.Stdin = nil
	cmd.Stderr = nil
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		//cancelFunc()
		log.Error(err)
		return
	}

	cmd.Start()

	log.Info("New FFMpeg service started")

	buf := make([]byte, 10485760)
	for {
		if ctx.Err() == context.DeadlineExceeded {
			videoServer.removeResource(resource)
			log.Error("Closing FFMpeg service.")
			return
		}
		n, err := stdout.Read(buf)
		if err != nil {
			// TODO: Check if someone else needs this FFMpeg service.
			//cancelFunc()
			if err == io.EOF {
				videoServer.removeResource(resource)
			}
			log.Error(err)
			return
		}
		outChan <- buf[:n]
	}

}
