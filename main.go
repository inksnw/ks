package main

import (
	"awesomeProject/api/resources"
	"awesomeProject/api/websocket"
	"awesomeProject/informer"
	"awesomeProject/k8sutils"
	"context"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/phuslu/log"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"net/url"
	"os"
	"strings"
	"time"

	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/cache"
)

type APIServer struct {
	Server          *http.Server
	container       *restful.Container
	InformerFactory informers.SharedInformerFactory
	RuntimeCache    cache.Cache
}

func main() {
	apiServer, err := NewAPIServer()
	if err != nil {
		panic(err)
	}
	apiServer.PrepareRun()
	apiServer.Run()
}

func (s *APIServer) waitForResourceSync() error {
	k8sGVRs := []schema.GroupVersionResource{
		{Group: "", Version: "v1", Resource: "namespaces"},
		{Group: "", Version: "v1", Resource: "nodes"},
		{Group: "", Version: "v1", Resource: "resourcequotas"},
		{Group: "", Version: "v1", Resource: "pods"},
		{Group: "", Version: "v1", Resource: "services"},
		{Group: "", Version: "v1", Resource: "persistentvolumeclaims"},
		{Group: "", Version: "v1", Resource: "persistentvolumes"},
		{Group: "", Version: "v1", Resource: "secrets"},
		{Group: "", Version: "v1", Resource: "configmaps"},
		{Group: "", Version: "v1", Resource: "serviceaccounts"},

		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "roles"},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "rolebindings"},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterroles"},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterrolebindings"},
		{Group: "apps", Version: "v1", Resource: "deployments"},
		{Group: "apps", Version: "v1", Resource: "daemonsets"},
		{Group: "apps", Version: "v1", Resource: "replicasets"},
		{Group: "apps", Version: "v1", Resource: "statefulsets"},
		{Group: "apps", Version: "v1", Resource: "controllerrevisions"},
		{Group: "storage.k8s.io", Version: "v1", Resource: "storageclasses"},
		{Group: "batch", Version: "v1", Resource: "jobs"},
		{Group: "batch", Version: "v1", Resource: "cronjobs"},
		{Group: "networking.k8s.io", Version: "v1", Resource: "ingresses"},
		{Group: "autoscaling", Version: "v2", Resource: "horizontalpodautoscalers"},
		{Group: "networking.k8s.io", Version: "v1", Resource: "networkpolicies"},
	}

	for _, gvr := range k8sGVRs {
		in, err := s.InformerFactory.ForResource(gvr)
		event := informer.CommonEvent{Resource: gvr.Resource}
		in.Informer().AddEventHandler(event)
		if err != nil {
			log.Error().Msgf("cannot create informer for %s", gvr)
			return err
		}
	}
	ctx := context.TODO()
	stopCh := ctx.Done()

	s.InformerFactory.Start(stopCh)
	s.InformerFactory.WaitForCacheSync(stopCh)

	go s.RuntimeCache.Start(ctx)
	s.RuntimeCache.WaitForCacheSync(ctx)

	log.Info().Msgf("Finished caching objects")

	return nil

}

func (s *APIServer) Run() {
	err := s.waitForResourceSync()
	if err != nil {
		panic(err)
	}
	log.Debug().Msgf("Start listening on %s", s.Server.Addr)
	err = s.Server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
func (s *APIServer) installKsAPIs() {
	resources.AddToContainer(s.container, s.InformerFactory, s.RuntimeCache)
}

func (s *APIServer) installWS() {
	websocket.AddToContainer(s.container)
}
func cors(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	chain.ProcessFilter(req, resp)

	allow := req.Request.Header.Get("Origin")
	resp.Header().Set("Access-Control-Allow-Origin", allow)
	resp.Header().Set("Access-Control-Allow-Credentials", "true")
	resp.Header().Set("Access-Control-Allow-Headers", "Content-Type,Access-Token")

}

func (s *APIServer) PrepareRun() {
	s.container = restful.NewContainer()

	s.container.Filter(cors)
	s.container.Router(restful.CurlyRouter{})
	s.installKsAPIs()
	s.installWS()

	for _, ws := range s.container.RegisteredWebServices() {
		log.Debug().Msgf("%s", ws.RootPath())
	}

	s.Server.Handler = s.container

	s.Server.Handler = WithKubeAPIServer(s.Server.Handler, k8sutils.K8sRestConfig(), &errorResponder{})

}

type errorResponder struct{}

func (e *errorResponder) Error(w http.ResponseWriter, req *http.Request, err error) {
	log.Error().Msgf(err.Error())
	responsewriters.InternalError(w, req, err)
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

		log.Debug().Msgf("请求Origin: %s", allow)
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

func NewAPIServer() (APIServer, error) {
	apiServer := APIServer{}
	server := &http.Server{
		Addr: fmt.Sprintf(":9090"),
	}

	apiServer.Server = server
	apiServer.InformerFactory = informers.NewSharedInformerFactory(k8sutils.InitClient(), 1*time.Minute)
	sch := scheme.Scheme
	var AddToSchemes runtime.SchemeBuilder
	err := AddToSchemes.AddToScheme(sch)
	if err != nil {
		return APIServer{}, err
	}

	apiServer.RuntimeCache, err = cache.New(k8sutils.K8sRestConfig(), cache.Options{Scheme: sch})
	if err != nil {
		return APIServer{}, err
	}

	return apiServer, nil
}

func init() {
	if !log.IsTerminal(os.Stderr.Fd()) {
		return
	}
	log.DefaultLogger = log.Logger{
		TimeFormat: "15:04:05",
		Caller:     1,
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    true,
			EndWithMessage: true,
		},
	}
}
