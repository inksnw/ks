/*
Copyright 2019 The KubeSphere Authors.

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

package networkpolicy

import (
	"awesomeProject/models/resources"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"

	"awesomeProject/api"
	"awesomeProject/apiserver/query"
)

type networkpolicyGetter struct {
	informers informers.SharedInformerFactory
}

func New(informers informers.SharedInformerFactory) resources.Interface {
	return &networkpolicyGetter{informers: informers}
}

func (n networkpolicyGetter) Get(namespace, name string) (runtime.Object, error) {
	return n.informers.Networking().V1().NetworkPolicies().Lister().NetworkPolicies(namespace).Get(name)
}

func (n networkpolicyGetter) List(namespace string, query *query.Query) (*api.ListResult, error) {
	nps, err := n.informers.Networking().V1().NetworkPolicies().Lister().NetworkPolicies(namespace).List(query.Selector())
	if err != nil {
		return nil, err
	}

	var result []runtime.Object
	for _, item := range nps {
		result = append(result, item)
	}

	return resources.DefaultList(result, query, n.compare, n.filter), nil
}

func (n networkpolicyGetter) filter(item runtime.Object, filter query.Filter) bool {
	np, ok := item.(*v1.NetworkPolicy)
	if !ok {
		return false
	}

	return resources.DefaultObjectMetaFilter(np.ObjectMeta, filter)
}

func (n networkpolicyGetter) compare(left runtime.Object, right runtime.Object, field query.Field) bool {
	leftNP, ok := left.(*v1.NetworkPolicy)
	if !ok {
		return false
	}

	rightNP, ok := right.(*v1.NetworkPolicy)
	if !ok {
		return true
	}
	return resources.DefaultObjectMetaCompare(leftNP.ObjectMeta, rightNP.ObjectMeta, field)
}
