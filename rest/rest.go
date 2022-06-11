package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ks/rest/handlers"
	"github.com/ks/rest/middlewares"
	"github.com/ks/rest/websocket"
)

func InitRestApi() {

	r := gin.New()
	//r.Use(middlewares.Recv())
	r.Use(middlewares.CorsMiddleware())
	r.GET("/ws", wsCore.ServeWs)
	r.GET("/webshell", wsCore.PodConnect)
	r.GET("/nodeshell", wsCore.NodeConnect)

	podRoot := fmt.Sprintf("/api/v1/")

	rootGroup := r.Group(podRoot)
	actionFactory(rootGroup)
	err := r.Run()
	if err != nil {
		return
	}

}

func actionFactory(rootGroup *gin.RouterGroup) {
	ac := []handlers.Action{
		handlers.Pod{},
		handlers.Deploy{},
		handlers.NameSpace{},
		handlers.Ingress{},
		handlers.Secret{},
		handlers.Node{},
		handlers.Role{},
	}
	for _, i := range ac {
		handlers.GenRouter(rootGroup, i)
	}
}
