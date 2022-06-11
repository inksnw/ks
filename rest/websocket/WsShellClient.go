package wsCore

import (
	"github.com/gorilla/websocket"
	"github.com/phuslu/log"
)

type WsShellClient struct {
	client *websocket.Conn
}

func NewWsShellClient(client *websocket.Conn) *WsShellClient {
	return &WsShellClient{client: client}
}
func (t *WsShellClient) Write(p []byte) (n int, err error) {
	err = t.client.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
func (t *WsShellClient) Read(p []byte) (n int, err error) {
	_, b, err := t.client.ReadMessage()
	if err != nil {
		return 0, err
	}
	log.Info().Msgf("接到wsshell信息: %s", string(b))
	return copy(p, b), nil
}
