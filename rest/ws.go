package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/phuslu/log"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll file for changes with this period.
	filePeriod = 10 * time.Second
)

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Error().Msgf("ws握手失败 %s", err)
		}
		return
	}

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
			continue
		}
	}

}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func BindWebSocketRouter(r *gin.Engine) {

	http.HandleFunc("/ws", serveWs)
	addr := ":9090"
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
