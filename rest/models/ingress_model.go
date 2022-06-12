package models

import (
	"context"
	"github.com/ks/k8sutils"
	"k8s.io/api/networking/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"strings"
)

type Ingress struct {
	Name       string `json:"name,omitempty"`
	NameSpace  string `json:"namespace"`
	CreateTime string `json:"create_time"`
	Host       string `json:"host"`
	Key        string `json:"key,omitempty"`
}

func (t Ingress) List(list *networkingv1.IngressList) (rv []Ingress) {
	for idx, i := range list.Items {
		rv = append(rv, Ingress{
			Name:       i.Name,
			NameSpace:  i.Namespace,
			Host:       i.Spec.Rules[0].Host,
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
	Name        string         `json:"name"`
	Namespace   string         `json:"namespace"`
	Host        string         `json:"host"`
	Annotations string         `json:"annotations"`
	Paths       []*IngressPath `json:"paths"`
}

//todo 解析标签
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

func (t IngressPost) Create() error {
	className := "nginx"
	pathType := v1.PathType("Prefix")
	var ingressRules []v1.IngressRule
	var rulePaths []v1.HTTPIngressPath
	for _, pathCfg := range t.Paths {
		PortInt, _ := strconv.ParseInt(pathCfg.Port, 10, 32)
		rule := v1.HTTPIngressPath{
			Path:     pathCfg.Path,
			PathType: &pathType,
			Backend: v1.IngressBackend{
				Service: &v1.IngressServiceBackend{},
			},
		}
		rule.Backend.Service.Name = pathCfg.SvcName
		rule.Backend.Service.Port.Number = int32(PortInt)
		rulePaths = append(rulePaths, rule)
	}
	rule := v1.IngressRule{
		Host: t.Host,
		IngressRuleValue: v1.IngressRuleValue{
			HTTP: &v1.HTTPIngressRuleValue{
				Paths: rulePaths,
			},
		},
	}
	ingressRules = append(ingressRules, rule)

	ingress := &v1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        t.Name,
			Namespace:   t.Namespace,
			Annotations: t.parseAnnotations(t.Annotations),
		},
		Spec: v1.IngressSpec{
			IngressClassName: &className,
			Rules:            ingressRules,
		},
	}

	_, err := k8sutils.Client.NetworkingV1().Ingresses(t.Namespace).
		Create(context.Background(), ingress, metav1.CreateOptions{})
	return err
}
