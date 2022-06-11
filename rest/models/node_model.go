package models

import corev1 "k8s.io/api/core/v1"

type Node struct {
	Name string `json:"name"`
	Key  int    `json:"key"`
}

func (t Node) List(list *corev1.NodeList) (rv []Node) {
	n := 0
	for _, item := range list.Items {
		n = n + 1
		rv = append(rv, Node{
			Name: item.Name,
			Key:  n,
		})
	}
	return rv

}
