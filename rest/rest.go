package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ks/rest/handlers"
	"github.com/ks/rest/middlewares"
)

func InitRestApi() {

	r := gin.New()
	//r.Use(middlewares.Recv())
	r.Use(middlewares.CorsMiddleware())
	podRoot := fmt.Sprintf("/api/v1/")

	go BindWebSocketRouter(r)

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
	}
	for _, i := range ac {
		handlers.GenRouter(rootGroup, i)
	}
}
