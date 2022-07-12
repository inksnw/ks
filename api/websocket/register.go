package websocket

import (
	"github.com/emicklei/go-restful"
)

func NewWebService() *restful.WebService {
	webservice := restful.WebService{}
	webservice.Path("/kapis/ws").Produces(restful.MIME_JSON)

	return &webservice
}

func AddToContainer(container *restful.Container) {

	webservice := NewWebService()
	cors := restful.CrossOriginResourceSharing{
		AllowedDomains: []string{"http://localhost:3000"},
		AllowedHeaders: []string{"Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With"},
		AllowedMethods: []string{"POST, OPTIONS, GET, PUT, DELETE"},
		CookiesAllowed: true,
		Container:      container,
	}
	container.Filter(cors.Filter)

	webservice.Route(webservice.GET("/ws").To(ServeWs))
	webservice.Route(webservice.GET("/webshell").To(PodConnect))
	webservice.Route(webservice.GET("/webshell").To(NodeConnect))
	container.Add(webservice)
}
