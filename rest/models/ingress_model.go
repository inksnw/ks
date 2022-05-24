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

// path 配置
type IngressPath struct {
	Path    string `json:"path"`
	SvcName string `json:"svc_name"`
	Port    string `json:"port"`
}

//提交Ingress 对象
type IngressPost struct {
	Name      string         `json:"name"`
	Namespace string         `json:"namespace"`
	Host      string         `json:"host"`
	Paths     []*IngressPath `json:"paths"`
}
