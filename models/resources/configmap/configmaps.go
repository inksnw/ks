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

package configmap

import (
	"awesomeProject/models/resources"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"

	"awesomeProject/api"
	"awesomeProject/apiserver/query"
)

type configmapsGetter struct {
	informer informers.SharedInformerFactory
}

func New(sharedInformers informers.SharedInformerFactory) resources.Interface {
	return &configmapsGetter{informer: sharedInformers}
}

func (d *configmapsGetter) Get(namespace, name string) (runtime.Object, error) {
	return d.informer.Core().V1().ConfigMaps().Lister().ConfigMaps(namespace).Get(name)
}

func (d *configmapsGetter) List(namespace string, query *query.Query) (*api.ListResult, error) {
	configmaps, err := d.informer.Core().V1().ConfigMaps().Lister().ConfigMaps(namespace).List(query.Selector())
	if err != nil {
		return nil, err
	}

	var result []runtime.Object
	for _, configmap := range configmaps {
		result = append(result, configmap)
	}

	return resources.DefaultList(result, query, d.compare, d.filter), nil
}

func (d *configmapsGetter) compare(left runtime.Object, right runtime.Object, field query.Field) bool {

	leftCM, ok := left.(*corev1.ConfigMap)
	if !ok {
		return false
	}

	rightCM, ok := right.(*corev1.ConfigMap)
	if !ok {
		return false
	}

	return resources.DefaultObjectMetaCompare(leftCM.ObjectMeta, rightCM.ObjectMeta, field)
}

func (d *configmapsGetter) filter(object runtime.Object, filter query.Filter) bool {
	configMap, ok := object.(*corev1.ConfigMap)
	if !ok {
		return false
	}

	return resources.DefaultObjectMetaFilter(configMap.ObjectMeta, filter)
}
