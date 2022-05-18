package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ks/k8sutils"
	"github.com/ks/rest/models"
	"github.com/phuslu/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Deploy struct {
	client *kubernetes.Clientset
}

func (d Deploy) List(c *gin.Context) {
	ns := c.Param("ns")
	list, err := k8sutils.Client.AppsV1().Deployments(ns).List(c, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	log.Debug().Msgf("总条数 %d", len(list.Items))
	podList := models.Deployment{}
	c.JSON(200, models.Result{
		Items: podList.List(list),
		Total: len(list.Items),
	})
}

func (d Deploy) Detail(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (d Deploy) Apply(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (d Deploy) Delete(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (d Deploy) GetResource() string {
	//TODO implement me
	return "deployments"
}
