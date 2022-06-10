package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Action interface {
	List(c *gin.Context)
	Detail(c *gin.Context)
	Apply(c *gin.Context)
	Delete(c *gin.Context)
	Exec(c *gin.Context)
	GetResource() string
}

func GenRouter(r *gin.RouterGroup, action Action) {
	r.GET(fmt.Sprintf("/%s", action.GetResource()), action.List)
	r.GET(fmt.Sprintf("/namespaces/:ns/%s", action.GetResource()), action.List)
	r.GET(fmt.Sprintf("/namespaces/:ns/%s/:name", action.GetResource()), action.Detail)
	r.POST(fmt.Sprintf("/namespaces/:ns/%s", action.GetResource()), action.Apply)
	r.POST(fmt.Sprintf("/%s", action.GetResource()), action.Apply)
	r.POST(fmt.Sprintf("/namespaces/:ns/%s/:name/exec", action.GetResource()), action.Exec)
	r.PUT(fmt.Sprintf("/namespaces/:ns/%s", action.GetResource()), action.Apply)
	r.DELETE(fmt.Sprintf("/namespaces/:ns/%s/:name", action.GetResource()), action.Delete)

}
