package models

import (
	rbacv1 "k8s.io/api/rbac/v1"
)

type Role struct {
	Name       string `json:"name"`
	NameSpace  string `json:"name_space"`
	CreateTime string `json:"create_time"`
	Key        int    `json:"key"`
}

func (t Role) List(list *rbacv1.RoleList) (rv []Role) {
	n := 0
	for _, item := range list.Items {
		n = n + 1
		rv = append(rv, Role{
			Name:       item.Name,
			NameSpace:  item.Namespace,
			CreateTime: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Key:        n,
		})
	}
	return rv

}
