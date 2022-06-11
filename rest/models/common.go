package models

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"strings"
)

type Result struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Items interface{} `json:"items"`
	Total int         `json:"total"`
}

func GetImagesByPod(containers []corev1.Container) string {
	images := containers[0].Image
	short := strings.Split(images, "@")[0]
	if imgLen := len(containers); imgLen > 1 {
		short += fmt.Sprintf("+其他%d个镜像", imgLen-1)
	}
	return short
}

func PosIsReady(pod *corev1.Pod) bool {
	if pod.Status.Phase != "Running" {
		return false
	}
	for _, condition := range pod.Status.Conditions {
		if condition.Status != "True" {
			return false
		}
	}
	for _, rg := range pod.Spec.ReadinessGates {
		for _, condition := range pod.Status.Conditions {
			if condition.Type == rg.ConditionType && condition.Status != "True" {
				return false
			}
		}
	}
	return true
}
