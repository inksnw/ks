package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ks/k8sutils"
	"github.com/ks/rest/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Node struct {
}

func (n Node) List(c *gin.Context) {
	list, err := k8sutils.Client.CoreV1().Nodes().List(c, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	nodeList := models.Node{}
	c.JSON(200,
		models.Result{
			Items: nodeList.List(list),
			Total: len(list.Items),
		},
	)
}

func (n Node) Detail(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (n Node) Apply(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (n Node) Delete(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (n Node) GetResource() string {
	//TODO implement me
	return "nodes"
}
