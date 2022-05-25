package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ks/k8sutils"
	"github.com/ks/rest/models"
	"github.com/phuslu/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Ingress struct {
}

func (i Ingress) List(c *gin.Context) {
	ns := c.Query("ns")
	list, err := k8sutils.Client.NetworkingV1().Ingresses(ns).List(c, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	log.Debug().Msgf("总条数 %d", len(list.Items))
	podList := models.Ingress{}
	c.JSON(200, models.Result{
		Items: podList.List(list),
		Total: len(list.Items),
	})
}

func (i Ingress) Detail(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (i Ingress) Apply(c *gin.Context) {

	postModel := models.IngressPost{}
	err := c.BindJSON(&postModel)
	if err != nil {
		panic(err)
	}
	postModel.Create()

	c.JSON(200, postModel)
}

func (i Ingress) Delete(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (i Ingress) GetResource() string {
	return "ingress"
}
