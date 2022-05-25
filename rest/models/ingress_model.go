package models

import (
	"context"
	"github.com/ks/k8sutils"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
	"strings"
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
	Port    int    `json:"port"`
}

//提交Ingress 对象
type IngressPost struct {
	Name        string         `json:"name"`
	Namespace   string         `json:"namespace"`
	Host        string         `json:"host"`
	Annotations string         `json:"annotations"`
	Paths       []*IngressPath `json:"paths"`
}

//解析标签
func (t IngressPost) parseAnnotations(annos string) map[string]string {
	replace := []string{"\t", " ", "\n", "\r\n"}
	for _, r := range replace {
		annos = strings.ReplaceAll(annos, r, "")
	}
	ret := make(map[string]string)
	list := strings.Split(annos, ";")
	for _, item := range list {
		annos := strings.Split(item, ":")
		if len(annos) == 2 {
			ret[annos[0]] = annos[1]
		}
	}
	return ret

}

func (t IngressPost) Create() {
	className := "nginx"
	var ingressRules []v1beta1.IngressRule
	var rulePaths []v1beta1.HTTPIngressPath
	for _, pathCfg := range t.Paths {
		rulePaths = append(rulePaths, v1beta1.HTTPIngressPath{
			Path: pathCfg.Path,
			Backend: v1beta1.IngressBackend{
				ServiceName: pathCfg.SvcName,
				ServicePort: intstr.FromInt(pathCfg.Port),
			},
		})
	}
	rule := v1beta1.IngressRule{
		Host: t.Host,
		IngressRuleValue: v1beta1.IngressRuleValue{
			HTTP: &v1beta1.HTTPIngressRuleValue{
				Paths: rulePaths,
			},
		},
	}
	ingressRules = append(ingressRules, rule)

	ingress := &v1beta1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        t.Name,
			Namespace:   t.Namespace,
			Annotations: t.parseAnnotations(t.Annotations),
		},
		Spec: v1beta1.IngressSpec{
			IngressClassName: &className,
			Rules:            ingressRules,
		},
	}

	_, err := k8sutils.Client.NetworkingV1beta1().Ingresses(t.Namespace).
		Create(context.Background(), ingress, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}
