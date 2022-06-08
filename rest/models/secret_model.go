package models

import (
	"context"
	"github.com/ks/k8sutils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Secret struct {
	Name      string `json:"name"`
	NameSpace string `json:"name_space"`
	Key       int    `json:"key"`
}

type SecretPost struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

func (p SecretPost) Create() {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.Name,
			Namespace: p.Namespace,
		},
	}
	secret.Type = corev1.SecretTypeOpaque
	secret.StringData = map[string]string{
		p.Key: p.Value,
	}

	_, err := k8sutils.Client.CoreV1().Secrets(p.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

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
