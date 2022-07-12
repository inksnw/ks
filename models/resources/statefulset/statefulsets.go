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

package statefulset

import (
	"awesomeProject/models/resources"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"

	"awesomeProject/api"
	"awesomeProject/apiserver/query"
)

const (
	statusStopped  = "stopped"
	statusRunning  = "running"
	statusUpdating = "updating"
)

type statefulSetGetter struct {
	sharedInformers informers.SharedInformerFactory
}

func New(sharedInformers informers.SharedInformerFactory) resources.Interface {
	return &statefulSetGetter{sharedInformers: sharedInformers}
}

func (d *statefulSetGetter) Get(namespace, name string) (runtime.Object, error) {
	return d.sharedInformers.Apps().V1().StatefulSets().Lister().StatefulSets(namespace).Get(name)
}

func (d *statefulSetGetter) List(namespace string, query *query.Query) (*api.ListResult, error) {
	// first retrieves all statefulSets within given namespace
	statefulSets, err := d.sharedInformers.Apps().V1().StatefulSets().Lister().StatefulSets(namespace).List(query.Selector())
	if err != nil {
		return nil, err
	}

	var result []runtime.Object
	for _, deployment := range statefulSets {
		result = append(result, deployment)
	}

	return resources.DefaultList(result, query, d.compare, d.filter), nil
}

func (d *statefulSetGetter) compare(left runtime.Object, right runtime.Object, field query.Field) bool {

	leftStatefulSet, ok := left.(*appsv1.StatefulSet)
	if !ok {
		return false
	}

	rightStatefulSet, ok := right.(*appsv1.StatefulSet)
	if !ok {
		return false
	}

	return resources.DefaultObjectMetaCompare(leftStatefulSet.ObjectMeta, rightStatefulSet.ObjectMeta, field)
}

func (d *statefulSetGetter) filter(object runtime.Object, filter query.Filter) bool {
	statefulSet, ok := object.(*appsv1.StatefulSet)
	if !ok {
		return false
	}

	switch filter.Field {
	case query.FieldStatus:
		return statefulSetStatus(statefulSet) == string(filter.Value)
	default:
		return resources.DefaultObjectMetaFilter(statefulSet.ObjectMeta, filter)
	}

}

func statefulSetStatus(item *appsv1.StatefulSet) string {
	if item.Spec.Replicas != nil {
		if item.Status.ReadyReplicas == 0 && *item.Spec.Replicas == 0 {
			return statusStopped
		} else if item.Status.ReadyReplicas == *item.Spec.Replicas {
			return statusRunning
		} else {
			return statusUpdating
		}
	}
	return statusStopped
}
