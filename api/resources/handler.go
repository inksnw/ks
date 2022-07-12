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

package resources

import (
	"awesomeProject/api"
	"awesomeProject/apiserver/query"
	"awesomeProject/models/resources/resource"
	"github.com/emicklei/go-restful"
	"github.com/phuslu/log"
)

type Handler struct {
	resourceGetter *resource.ResourceGetter
}

func New(resourceGetterV1alpha3 *resource.ResourceGetter) *Handler {
	return &Handler{
		resourceGetter: resourceGetterV1alpha3,
	}
}

func (h *Handler) handleGetResources(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	resourceType := request.PathParameter("resources")
	name := request.PathParameter("name")

	// use informers to retrieve resources
	result, err := h.resourceGetter.Get(resourceType, namespace, name)
	if err == nil {
		response.WriteEntity(result)
		return
	}

	if err != resource.ErrResourceNotSupported {
		log.Error().Msgf("%s %s", err, resourceType)
		api.HandleInternalError(response, nil, err)
		return
	}

}

// handleListResources retrieves resources
func (h *Handler) handleListResources(request *restful.Request, response *restful.Response) {
	queryArg := query.ParseQueryParameter(request)
	resourceType := request.PathParameter("resources")
	namespace := request.PathParameter("namespace")

	result, err := h.resourceGetter.List(resourceType, namespace, queryArg)
	if err == nil {
		response.WriteEntity(result)
		return
	}

	if err != resource.ErrResourceNotSupported {
		log.Error().Msgf("%s %s", err, resourceType)
		api.HandleInternalError(response, request, err)
		return
	}
	response.WriteEntity(result)
}
