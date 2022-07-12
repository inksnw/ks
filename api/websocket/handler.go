/*
Copyright 2020 KubeSphere Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package websocket

import (
	"awesomeProject/helpers"
	"awesomeProject/k8sutils"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/gorilla/websocket"
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

func ServeWs(req *restful.Request, resp *restful.Response) {

	ws, err := UpGrader.Upgrade(resp.ResponseWriter, req.Request, nil)
	if err != nil {
		log.Error().Msgf("ws握手失败 %s", err)
		return
	}
	log.Info().Msgf("新建ws连接来自 %s", ws.RemoteAddr())
	ClientMap.Store(ws)
}
func PodConnect(req *restful.Request, resp *restful.Response) {
	wsClient, err := UpGrader.Upgrade(resp.ResponseWriter, req.Request, nil)
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

func NodeConnect(req *restful.Request, resp *restful.Response) {
	wsClient, err := UpGrader.Upgrade(resp.ResponseWriter, req.Request, nil)
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
