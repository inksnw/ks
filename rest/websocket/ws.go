package wsCore

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ks/k8sutils"
	"github.com/ks/rest/helpers"
	"github.com/phuslu/log"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
)

var (
	UpGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func ServeWs(c *gin.Context) {
	ws, err := UpGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error().Msgf("ws握手失败 %s", err)
		return
	}
	log.Info().Msgf("新建ws连接来自 %s", ws.RemoteAddr())
	ClientMap.Store(ws)
}
func PodConnect(c *gin.Context) {
	wsClient, err := UpGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}
	shellClient := NewWsShellClient(wsClient)

	cli := k8sutils.InitClient()
	err = helpers.HandleCommand(cli, k8sutils.K8sRestConfig(), []string{"sh"}).
		Stream(remotecommand.StreamOptions{
			Stdin:  shellClient,
			Stdout: shellClient,
			Stderr: shellClient,
			Tty:    true,
		})
	return
}

func NodeConnect(c *gin.Context) {
	wsClient, err := UpGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}
	shellClient := NewWsShellClient(wsClient)
	session, err := helpers.SSHConnect(helpers.TempSSHUser, helpers.TempSSHPWD, helpers.TempSSHIP, 22)
	fmt.Println(err)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.Stdout = shellClient
	session.Stderr = shellClient
	session.Stdin = shellClient
	err = session.RequestPty("xterm-256color", 300, 500, helpers.NodeShellModes)
	if err != nil {
		panic(err)
	}

	err = session.Run("sh")
	if err != nil {
		panic(err)
	}
	return
}
