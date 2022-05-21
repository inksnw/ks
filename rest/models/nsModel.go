package models

import (
	corev1 "k8s.io/api/core/v1"
	"strconv"
)

type NameSpace struct {
	Name string `json:"name,omitempty"`
	Key  string `json:"key,omitempty"`
}

func (t NameSpace) List(list *corev1.NamespaceList) (rv []NameSpace) {
	for idx, i := range list.Items {
		rv = append(rv, NameSpace{
			Name: i.Name,
			Key:  strconv.Itoa(idx),
		})
	}
	return rv
}
