package resources

import (
	"awesomeProject/api"
	"awesomeProject/apiserver/query"
	"awesomeProject/models/resources/resource"
	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/cache"
)

const (
	GroupName = "resources.kubesphere.io"
	ok        = "OK"
)

func NewWebService(gv schema.GroupVersion) *restful.WebService {
	webservice := restful.WebService{}
	webservice.Path("/kapis/" + gv.String()).
		Produces(restful.MIME_JSON)

	return &webservice
}

var GroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha3"}

func AddToContainer(c *restful.Container, informerFactory informers.SharedInformerFactory, cache cache.Cache) {

	webservice := NewWebService(GroupVersion)
	handler := New(resource.NewResourceGetter(informerFactory, cache))

	webservice.Route(webservice.GET("/{resources}").
		To(handler.handleListResources).
		Doc("Cluster level resources").
		Param(webservice.PathParameter("resources", "cluster level resource type, e.g. pods,jobs,configmaps,services.")).
		Param(webservice.QueryParameter(query.ParameterName, "name used to do filtering").Required(false)).
		Param(webservice.QueryParameter(query.ParameterPage, "page").Required(false).DataFormat("page=%d").DefaultValue("page=1")).
		Param(webservice.QueryParameter(query.ParameterLimit, "limit").Required(false)).
		Param(webservice.QueryParameter(query.ParameterAscending, "sort parameters, e.g. reverse=true").Required(false).DefaultValue("ascending=false")).
		Param(webservice.QueryParameter(query.ParameterOrderBy, "sort parameters, e.g. orderBy=createTime")).
		Returns(http.StatusOK, ok, api.ListResult{}))

	webservice.Route(webservice.GET("/{resources}/{name}").
		To(handler.handleGetResources).
		Doc("Cluster level resource").
		Param(webservice.PathParameter("resources", "cluster level resource type, e.g. pods,jobs,configmaps,services.")).
		Param(webservice.PathParameter("name", "the name of the clustered resources")).
		Returns(http.StatusOK, api.StatusOK, nil))

	webservice.Route(webservice.GET("/namespaces/{namespace}/{resources}").
		To(handler.handleListResources).
		Doc("Namespace level resource query").
		Param(webservice.PathParameter("namespace", "the name of the project")).
		Param(webservice.PathParameter("resources", "namespace level resource type, e.g. pods,jobs,configmaps,services.")).
		Param(webservice.QueryParameter(query.ParameterName, "name used to do filtering").Required(false)).
		Param(webservice.QueryParameter(query.ParameterPage, "page").Required(false).DataFormat("page=%d").DefaultValue("page=1")).
		Param(webservice.QueryParameter(query.ParameterLimit, "limit").Required(false)).
		Param(webservice.QueryParameter(query.ParameterAscending, "sort parameters, e.g. reverse=true").Required(false).DefaultValue("ascending=false")).
		Param(webservice.QueryParameter(query.ParameterOrderBy, "sort parameters, e.g. orderBy=createTime")).
		Returns(http.StatusOK, ok, api.ListResult{}))

	webservice.Route(webservice.GET("/namespaces/{namespace}/{resources}/{name}").
		To(handler.handleGetResources).
		Doc("Namespace level get resource query").
		Param(webservice.PathParameter("namespace", "the name of the project")).
		Param(webservice.PathParameter("resources", "namespace level resource type, e.g. pods,jobs,configmaps,services.")).
		Param(webservice.PathParameter("name", "the name of resource")).
		Returns(http.StatusOK, ok, api.ListResult{}))
	c.Add(webservice)
}
