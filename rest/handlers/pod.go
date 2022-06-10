package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ks/k8sutils"
	"github.com/ks/rest/models"
	"github.com/phuslu/log"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
	"os"
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
	log.Debug().Msgf("总条数 %d", len(list.Items))
	podList := models.Pod{}
	c.JSON(200, models.Result{
		Items: podList.List(list),
		Total: len(list.Items),
	})
}

func (p Pod) Detail(c *gin.Context) {
	ns := c.Param("ns")
	name := c.Param("name")
	req := k8sutils.Client.CoreV1().Pods(ns).GetLogs(name, &corev1.PodLogOptions{Follow: true})
	reader, err := req.Stream(context.Background())
	if err != nil {
		panic(err)
	}

	for {
		buf := make([]byte, 1024)
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			break
		}
		if n > 0 {
			c.Writer.Write([]byte(string(buf[0:n])))
			c.Writer.(http.Flusher).Flush()
		}
	}
}

func (p Pod) Exec(c *gin.Context) {
	ns := c.Param("ns")
	name := c.Param("name")
	cmd := []string{"sh", "-c", "ls"}

	option := &corev1.PodExecOptions{
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
		Container: "mynginx",
		Command:   cmd,
	}
	req := k8sutils.Client.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(ns).Name(name).
		SubResource("exec").
		VersionedParams(
			option,
			scheme.ParameterCodec,
		)
	exec, err := remotecommand.NewSPDYExecutor(
		k8sutils.K8sRestConfig(),
		"POST",
		req.URL())
	if err != nil {
		panic(err)
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: c.Writer,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		panic(err)
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
