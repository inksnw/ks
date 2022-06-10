package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ks/k8sutils"
	"github.com/ks/rest/models"
	"github.com/phuslu/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NameSpace struct {
	client *kubernetes.Clientset
}

func (n NameSpace) Exec(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (n NameSpace) List(c *gin.Context) {
	list, err := k8sutils.Client.CoreV1().Namespaces().List(c, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	log.Debug().Msgf("总条数 %d", len(list.Items))
	podList := models.NameSpace{}
	c.JSON(200, models.Result{
		Items: podList.List(list),
		Total: len(list.Items),
	})
}

func (n NameSpace) Detail(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (n NameSpace) Apply(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (n NameSpace) Delete(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (n NameSpace) GetResource() string {
	return "namespaces"
}
