package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ks/k8sutils"
	"github.com/phuslu/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Pod struct {
}

func (p Pod) GetResource() string {
	//TODO implement me
	return "pods"
}

type result struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Items interface{} `json:"items"`
	Total int         `json:"total"`
}

type podInfo struct {
	Name       string      `json:"name"`
	NameSpace  string      `json:"name_space"`
	CreateTime metav1.Time `json:"create_time"`
	Status     string      `json:"status"`
	Key        int         `json:"key"`
}

func (p Pod) List(c *gin.Context) {
	client := k8sutils.InitClient()
	list, err := client.CoreV1().Pods("").List(c, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	log.Debug().Msgf("总条数 %d", len(list.Items))
	var rv []podInfo
	var n = 0
	for _, i := range list.Items {
		n = n + 1
		rv = append(rv, podInfo{
			Key:        n,
			Name:       i.Name,
			NameSpace:  i.Namespace,
			CreateTime: i.CreationTimestamp,
			Status:     string(i.Status.Phase),
		})
	}
	c.JSON(200, result{
		Items: rv,
		Total: len(list.Items),
	})
}

func (p Pod) Detail(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (p Pod) Apply(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (p Pod) Delete(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
