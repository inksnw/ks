/*
Copyright 2020 The KubeSphere Authors.

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

package resource

import (
	"awesomeProject/api"
	"awesomeProject/apiserver/query"
	"awesomeProject/models/resources/configmap"
	"awesomeProject/models/resources/daemonset"
	"awesomeProject/models/resources/deployment"
	"awesomeProject/models/resources/ingress"
	"awesomeProject/models/resources/job"
	"awesomeProject/models/resources/namespace"
	"awesomeProject/models/resources/networkpolicy"
	"awesomeProject/models/resources/node"
	"awesomeProject/models/resources/persistentvolume"
	"awesomeProject/models/resources/pod"
	"awesomeProject/models/resources/role"
	"awesomeProject/models/resources/rolebinding"
	"awesomeProject/models/resources/secret"
	"awesomeProject/models/resources/service"
	"awesomeProject/models/resources/serviceaccount"
	"awesomeProject/models/resources/statefulset"
	"errors"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"sigs.k8s.io/controller-runtime/pkg/cache"

	"awesomeProject/models/resources"
	"k8s.io/client-go/informers"
)

var ErrResourceNotSupported = errors.New("resource is not supported")

type ResourceGetter struct {
	clusterResourceGetters    map[schema.GroupVersionResource]resources.Interface
	namespacedResourceGetters map[schema.GroupVersionResource]resources.Interface
}

func NewResourceGetter(factory informers.SharedInformerFactory, cache cache.Cache) *ResourceGetter {
	namespacedResourceGetters := make(map[schema.GroupVersionResource]resources.Interface)
	clusterResourceGetters := make(map[schema.GroupVersionResource]resources.Interface)

	namespacedResourceGetters[schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}] = deployment.New(factory)
	namespacedResourceGetters[schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "daemonsets"}] = daemonset.New(factory)
	namespacedResourceGetters[schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"}] = statefulset.New(factory)
	namespacedResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}] = service.New(factory)
	namespacedResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}] = configmap.New(factory)
	namespacedResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "secrets"}] = secret.New(factory)
	namespacedResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}] = pod.New(factory)
	namespacedResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "serviceaccounts"}] = serviceaccount.New(factory)
	namespacedResourceGetters[schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "ingresses"}] = ingress.New(factory)
	namespacedResourceGetters[schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "networkpolicies"}] = networkpolicy.New(factory)
	namespacedResourceGetters[schema.GroupVersionResource{Group: "batch", Version: "v1", Resource: "jobs"}] = job.New(factory)
	clusterResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "persistentvolumes"}] = persistentvolume.New(factory)
	namespacedResourceGetters[rbacv1.SchemeGroupVersion.WithResource("rolebindings")] = rolebinding.New(factory)
	namespacedResourceGetters[rbacv1.SchemeGroupVersion.WithResource("roles")] = role.New(factory)
	clusterResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "nodes"}] = node.New(factory)
	clusterResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"}] = namespace.New(factory)

	return &ResourceGetter{
		namespacedResourceGetters: namespacedResourceGetters,
		clusterResourceGetters:    clusterResourceGetters,
	}
}

// TryResource will retrieve a getter with resource name, it doesn't guarantee find resource with correct group version
// need to refactor this use schema.GroupVersionResource
func (r *ResourceGetter) TryResource(clusterScope bool, resource string) resources.Interface {
	if clusterScope {
		for k, v := range r.clusterResourceGetters {
			if k.Resource == resource {
				return v
			}
		}
	}
	for k, v := range r.namespacedResourceGetters {
		if k.Resource == resource {
			return v
		}
	}
	return nil
}

func (r *ResourceGetter) Get(resource, namespace, name string) (runtime.Object, error) {
	clusterScope := namespace == ""
	getter := r.TryResource(clusterScope, resource)
	if getter == nil {
		return nil, ErrResourceNotSupported
	}
	return getter.Get(namespace, name)
}

func (r *ResourceGetter) List(resource, namespace string, query *query.Query) (*api.ListResult, error) {
	clusterScope := namespace == ""
	getter := r.TryResource(clusterScope, resource)
	if getter == nil {
		return nil, ErrResourceNotSupported
	}
	return getter.List(namespace, query)
}
