package models

import (
	networkingv1 "k8s.io/api/networking/v1"
	"strconv"
)

type Ingress struct {
	Name       string `json:"name,omitempty"`
	NameSpace  string `json:"name_space"`
	CreateTime string `json:"create_time"`
	Key        string `json:"key,omitempty"`
}

func (t Ingress) List(list *networkingv1.IngressList) (rv []Ingress) {
	for idx, i := range list.Items {
		rv = append(rv, Ingress{
			Name:       i.Name,
			NameSpace:  i.Namespace,
			Key:        strconv.Itoa(idx),
			CreateTime: i.CreationTimestamp.Format("2006-01-02 15:04:05"),
		})
	}
	return rv
}
