package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/phuslu/log"
	"net/http"
	"time"
)

var (
	upGrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func ServeWs(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error().Msgf("ws握手失败 %s", err)
		return
	}
	log.Info().Msgf("新建ws连接来自 %s", ws.RemoteAddr())

	go writer(ws)
	reader(ws)
}

func writer(ws *websocket.Conn) {

	for {
		time.Sleep(5 * time.Second)
		data := make(map[string]string)
		data["hh"] = "sss"
		err := ws.WriteJSON(data)
		if err != nil {
			log.Error().Msgf("发送失败 %s", err)
			break
		}
	}

}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}
