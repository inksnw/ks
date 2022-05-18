package models

import (
	appsv1 "k8s.io/api/apps/v1"
)

type Deployment struct {
	Name       string   `json:"name,omitempty"`
	NameSpace  string   `json:"name_space,omitempty"`
	Replicas   [3]int32 `json:"replicas,omitempty"` //3个值，分别是总副本数，可用副本数 ，不可用副本数
	Images     string   `json:"images,omitempty"`
	IsComplete bool     `json:"is_complete,omitempty"` //是否完成
	Message    string   `json:"message,omitempty"`     // 显示错误信息
	CreateTime string   `json:"create_time,omitempty"`
	Pods       []*Pod   `json:"pods,omitempty"`
}

func (t Deployment) List(list *appsv1.DeploymentList) (rv []Deployment) {
	n := 0
	for _, i := range list.Items {
		n = n + 1
		rv = append(rv, Deployment{
			Name:       i.Name,
			NameSpace:  i.Namespace,
			Replicas:   [3]int32{},
			Images:     "",
			IsComplete: false,
			Message:    "",
			CreateTime: "",
			Pods:       nil,
		})
	}
	return rv

}
