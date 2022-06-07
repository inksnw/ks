package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ks/k8sutils"
	"github.com/ks/rest/models"
	"github.com/phuslu/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Secret struct {
}

func (s Secret) List(c *gin.Context) {
	ns := c.Query("ns")

	list, err := k8sutils.Client.CoreV1().Secrets(ns).List(c, metav1.ListOptions{})

	if err != nil {
		panic(err)
	}
	log.Debug().Msgf("总条数 %d", len(list.Items))
	secretList := models.Secret{}
	c.JSON(200, models.Result{
		Items: secretList.List(list),
		Total: len(list.Items),
	})
}

func (s Secret) Detail(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (s Secret) Apply(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (s Secret) Delete(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (s Secret) GetResource() string {
	//TODO implement me
	return "secret"
}