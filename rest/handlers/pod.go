package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ks/k8sutils"
	"github.com/ks/rest/models"
	"github.com/phuslu/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Pod struct {
	client *kubernetes.Clientset
}

func (p Pod) SetClient(cli *kubernetes.Clientset) {
	p.client = cli
}

func (p Pod) GetResource() string {
	//TODO implement me
	return "pods"
}

func (p Pod) List(c *gin.Context) {

	list, err := k8sutils.Client.CoreV1().Pods("").List(c, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	log.Debug().Msgf("总条数 %d", len(list.Items))
	podList := models.Pod{}
	c.JSON(200, models.Result{
		Items: podList.List(list),
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
