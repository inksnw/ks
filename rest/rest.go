package rest

import (
	"bytes"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/ks/k8sutils"
	"github.com/ks/rest/websocket"
	"github.com/phuslu/log"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"net/http"
	"net/url"
	rt "runtime"
	"strings"
)

func NewWebService() *restful.WebService {
	webservice := restful.WebService{}
	webservice.Path("/kapis/").Produces(restful.MIME_JSON)

	return &webservice
}

func logStackOnRecover(panicReason interface{}, w http.ResponseWriter) {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("recover from panic situation: - %v\r\n", panicReason))
	for i := 2; ; i += 1 {
		_, file, line, ok := rt.Caller(i)
		if !ok {
			break
		}
		buffer.WriteString(fmt.Sprintf("    %s:%d\r\n", file, line))
	}
	klog.Errorln(buffer.String())

	headers := http.Header{}
	if ct := w.Header().Get("Content-Type"); len(ct) > 0 {
		headers.Set("Accept", ct)
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal server error"))
}
func InitRestApi() {

	container := restful.NewContainer()
	container.RecoverHandler(func(panicReason interface{}, httpWriter http.ResponseWriter) {
		logStackOnRecover(panicReason, httpWriter)
	})
	//container.ServeMux
	cors := restful.CrossOriginResourceSharing{
		AllowedDomains: []string{"http://localhost:3000"},
		AllowedHeaders: []string{"Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With"},
		AllowedMethods: []string{"POST, OPTIONS, GET, PUT, DELETE"},
		CookiesAllowed: true,
		Container:      container,
	}
	container.Filter(cors.Filter)

	ws := NewWebService()
	ws.Route(ws.GET("/ws").To(wsCore.ServeWs))
	ws.Route(ws.GET("/webshell").To(wsCore.PodConnect))
	ws.Route(ws.GET("/webshell").To(wsCore.NodeConnect))
	container.Add(ws)

	server := &http.Server{
		Addr:    ":8080",
		Handler: container,
	}
	server.Handler = WithKubeAPIServer(server.Handler, k8sutils.K8sRestConfig(), &errorResponder{})

	for _, ins := range container.RegisteredWebServices() {
		log.Info().Msgf("服务地址: %s", ins.RootPath())
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
func WithKubeAPIServer(handler http.Handler, config *rest.Config, failed proxy.ErrorResponder) http.Handler {
	defaultTransport, err := rest.TransportFor(config)
	if err != nil {
		klog.Errorf("Unable to create transport from rest.Config: %v", err)
		return handler
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		if !IsKubernetesRequest(req) {
			handler.ServeHTTP(w, req)
			return
		}
		allow := req.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", allow)

		httpProxy := proxy.NewUpgradeAwareHandler(getLocation(req), defaultTransport, true, false, failed)
		httpProxy.UpgradeTransport = proxy.NewUpgradeRequestRoundTripper(defaultTransport, defaultTransport)
		httpProxy.ServeHTTP(w, req)
	})
}

func IsKubernetesRequest(req *http.Request) bool {
	currentParts := splitPath(req.URL.Path)
	if len(currentParts) == 0 {
		return true
	}
	APIPrefixes := sets.NewString("api", "apis")
	if APIPrefixes.Has(currentParts[0]) {
		return true
	}
	return false
}
func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return []string{}
	}
	return strings.Split(path, "/")
}

func getLocation(req *http.Request) *url.URL {
	config := k8sutils.K8sRestConfig()
	k8sUrl, _ := url.Parse(config.Host)
	req.URL.Host = k8sUrl.Host
	req.URL.Scheme = k8sUrl.Scheme
	return req.URL
}

type errorResponder struct{}

func (e *errorResponder) Error(w http.ResponseWriter, req *http.Request, err error) {
	klog.Error(err)
	responsewriters.InternalError(w, req, err)
}
