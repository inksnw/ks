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

package role

import (
	"awesomeProject/api"
	"awesomeProject/apiserver/query"
	"awesomeProject/models/resources"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"
)

type rolesGetter struct {
	sharedInformers informers.SharedInformerFactory
}

func New(sharedInformers informers.SharedInformerFactory) resources.Interface {
	return &rolesGetter{sharedInformers: sharedInformers}
}

func (d *rolesGetter) Get(namespace, name string) (runtime.Object, error) {
	return d.sharedInformers.Rbac().V1().Roles().Lister().Roles(namespace).Get(name)
}

func (d *rolesGetter) List(namespace string, query *query.Query) (*api.ListResult, error) {

	var roles []*rbacv1.Role
	var err error

	roles, err = d.sharedInformers.Rbac().V1().Roles().Lister().Roles(namespace).List(query.Selector())

	if err != nil {
		return nil, err
	}

	var result []runtime.Object
	for _, role := range roles {
		result = append(result, role)
	}

	return resources.DefaultList(result, query, d.compare, d.filter), nil
}

func (d *rolesGetter) compare(left runtime.Object, right runtime.Object, field query.Field) bool {

	leftRole, ok := left.(*rbacv1.Role)
	if !ok {
		return false
	}

	rightRole, ok := right.(*rbacv1.Role)
	if !ok {
		return false
	}

	return resources.DefaultObjectMetaCompare(leftRole.ObjectMeta, rightRole.ObjectMeta, field)
}

func (d *rolesGetter) filter(object runtime.Object, filter query.Filter) bool {
	role, ok := object.(*rbacv1.Role)

	if !ok {
		return false
	}

	return resources.DefaultObjectMetaFilter(role.ObjectMeta, filter)
}

func (d *rolesGetter) fetchAggregationRoles(namespace, name string) ([]*rbacv1.Role, error) {
	roles := make([]*rbacv1.Role, 0)
	_, err := d.Get(namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			return roles, nil
		}
		return nil, err
	}
	return roles, nil
}
