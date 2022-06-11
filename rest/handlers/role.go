package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ks/k8sutils"
	"github.com/ks/rest/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Role struct {
}

func (r Role) List(c *gin.Context) {
	ns := c.Query("ns")
	list, err := k8sutils.Client.RbacV1().Roles(ns).List(c, metav1.ListOptions{})
	role := models.Role{}

	if err != nil {
		panic(err)
	}
	c.JSON(200, models.Result{
		Items: role.List(list),
		Total: len(list.Items),
	})
}

func (r Role) Detail(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r Role) Apply(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r Role) Delete(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r Role) GetResource() string {
	return "role"
}
