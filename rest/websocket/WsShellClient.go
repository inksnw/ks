package wsCore

import (
	"github.com/gorilla/websocket"
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
	return copy(p, string(b)+"\n"), nil
}
