package models

import (
	corev1 "k8s.io/api/core/v1"
)

type Pod struct {
	Name       string   `json:"name"`
	NameSpace  string   `json:"name_space"`
	Images     string   `json:"images"`
	NodeName   string   `json:"node_name"`
	IP         []string `json:"IP"` // 第一个是 POD IP 第二个是 node ip
	CreateTime string   `json:"create_time"`
	Phase      string   `json:"phase"`    // pod 当前所处的阶段
	IsReady    bool     `json:"is_ready"` //判断pod 是否就绪
	Key        int      `json:"key"`
}

func (t Pod) List(list *corev1.PodList) (rv []Pod) {
	n := 0
	for _, pod := range list.Items {
		n = n + 1
		rv = append(rv, Pod{
			Name:       pod.Name,
			NameSpace:  pod.Namespace,
			Images:     GetImagesByPod(pod.Spec.Containers),
			NodeName:   pod.Spec.NodeName,
			IP:         []string{pod.Status.PodIP, pod.Status.HostIP},
			CreateTime: pod.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Phase:      string(pod.Status.Phase),
			IsReady:    PosIsReady(&pod),
			Key:        n,
		})
	}
	return rv

}
