package handlers

import (
	"context"
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
	ns := c.Param("ns")
	name := c.Param("name")
	ingress, err := k8sutils.Client.NetworkingV1beta1().Ingresses(ns).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	c.JSON(200, ingress)
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

	ns := c.Param("ns")
	name := c.Param("name")
	err := k8sutils.Client.NetworkingV1().Ingresses(ns).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
	c.JSON(200, models.Result{
		Msg: "删除成功",
	})
}

func (i Ingress) GetResource() string {
	return "ingress"
}
