package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ks/k8sutils"
	"github.com/ks/rest/models"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

type Pod struct {
}

func (p Pod) GetResource() string {
	//TODO implement me
	return "pods"
}

func (p Pod) List(c *gin.Context) {
	ns := c.Query("ns")

	list, err := k8sutils.Client.CoreV1().Pods(ns).List(c, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	podList := models.Pod{}
	c.JSON(200, models.Result{
		Items: podList.List(list),
		Total: len(list.Items),
	})
}

func (p Pod) Detail(c *gin.Context) {
	ns := c.Param("ns")
	name := c.Param("name")
	cc, _ := context.WithTimeout(c, time.Minute*30) //设置半小时超时时间。否则会造成内存泄露
	var tailLine int64 = 200
	req := k8sutils.Client.CoreV1().Pods(ns).GetLogs(name,
		&corev1.PodLogOptions{Follow: true,
			TailLines: &tailLine,
		},
	)
	reader, err := req.Stream(cc)
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	for {
		buf := make([]byte, 1024)
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			break
		}
		if n > 0 {
			c.Writer.Write(buf[0:n])
			c.Writer.(http.Flusher).Flush()
		}
	}
}

func (p Pod) Apply(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (p Pod) Delete(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
