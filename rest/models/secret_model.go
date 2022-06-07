package models

import (
	corev1 "k8s.io/api/core/v1"
)

type Secret struct {
	Name      string `json:"name"`
	NameSpace string `json:"name_space"`
	Key       int    `json:"key"`
}

func (t Secret) List(list *corev1.SecretList) (rv []Secret) {
	n := 0
	for _, item := range list.Items {
		n = n + 1
		rv = append(rv, Secret{
			Name:      item.Name,
			NameSpace: item.Namespace,
			Key:       n,
		})
	}
	return rv
}
