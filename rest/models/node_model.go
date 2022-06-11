package models

import (
	"github.com/ks/k8sutils"
	"github.com/ks/rest/helpers"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

type Node struct {
	Name string  `json:"name"`
	Cpu  float64 `json:"cpu"`
	Mem  float64 `json:"mem"`
	Key  int     `json:"key"`
}

func (t Node) List(list *corev1.NodeList) (rv []Node) {
	n := 0
	clientSet, _ := versioned.NewForConfig(k8sutils.K8sRestConfig())
	for _, item := range list.Items {
		n = n + 1
		cpu, mem := helpers.GetNodeUsage(clientSet, &item)
		rv = append(rv, Node{
			Name: item.Name,
			Cpu:  cpu,
			Mem:  mem,
			Key:  n,
		})
	}
	return rv

}
