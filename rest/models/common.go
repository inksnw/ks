package models

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type Result struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Items interface{} `json:"items"`
	Total int         `json:"total"`
}

func GetImages(dep appsv1.Deployment) string {
	return GetImagesByPod(dep.Spec.Template.Spec.Containers)
}
func GetImagesByPod(containers []corev1.Container) string {
	images := containers[0].Image
	if imgLen := len(containers); imgLen > 1 {
		images += fmt.Sprintf("+其他%d个镜像", imgLen-1)
	}
	return images
}
