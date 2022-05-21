package models

import (
	appsv1 "k8s.io/api/apps/v1"
	"strconv"
)

type Deployment struct {
	Name       string   `json:"name,omitempty"`
	NameSpace  string   `json:"name_space,omitempty"`
	Replicas   [3]int32 `json:"replicas,omitempty"` //3个值，分别是总副本数，可用副本数 ，不可用副本数
	Images     string   `json:"images,omitempty"`
	IsComplete bool     `json:"is_complete,omitempty"` //是否完成
	Message    string   `json:"message,omitempty"`     // 显示错误信息
	CreateTime string   `json:"create_time,omitempty"`
	Key        string   `json:"key,omitempty"`
}

func (t Deployment) List(list *appsv1.DeploymentList) (rv []Deployment) {
	for idx, i := range list.Items {
		rv = append(rv, Deployment{
			Name:       i.Name,
			NameSpace:  i.Namespace,
			Replicas:   [3]int32{i.Status.Replicas, i.Status.AvailableReplicas, i.Status.UnavailableReplicas},
			Images:     GetImagesByPod(i.Spec.Template.Spec.Containers),
			IsComplete: t.getDeploymentIsComplete(i),
			Message:    t.getDeploymentCondition(i),
			CreateTime: i.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Key:        strconv.Itoa(idx),
		})
	}
	return rv

}

func (t Deployment) getDeploymentIsComplete(dep appsv1.Deployment) bool {
	return dep.Status.Replicas == dep.Status.AvailableReplicas
}

func (t Deployment) getDeploymentCondition(dep appsv1.Deployment) string {
	for _, item := range dep.Status.Conditions {
		if string(item.Type) == "Available" && string(item.Status) != "True" {
			return item.Message
		}
	}
	return ""
}
