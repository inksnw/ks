package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ks/rest/handlers"
	"github.com/ks/rest/middlewares"
	"github.com/ks/rest/websocket"
)

func InitRestApi() {

	r := gin.New()
	//r.Use(middlewares.Recv())
	r.Use(middlewares.CorsMiddleware())
	r.GET("/ws", wsCore.ServeWs)

	r.GET("/send", wstest)
	podRoot := fmt.Sprintf("/api/v1/")

	rootGroup := r.Group(podRoot)
	actionFactory(rootGroup)
	err := r.Run()
	if err != nil {
		return
	}

}

func wstest(c *gin.Context) {
	var data = `
{
    "code": 0,
    "msg": "",
    "items": [
        {
            "name": "我是ws",
            "name_space": "ingress-nginx",
            "replicas": [
                1,
                1,
                0
            ],
            "images": "k8s.gcr.io/ingress-nginx/controller:v1.2.0@sha256:d8196e3bc1e72547c5dec66d6556c0ff92a23f6d0919b206be170bc90d5f9185",
            "is_complete": true,
            "key": "0"
        },
        {
            "name": "ingressgateway",
            "name_space": "istio-system",
            "replicas": [
                1,
                1,
                0
            ],
            "images": "auto",
            "is_complete": true,
            "key": "1"
        },
        {
            "name": "istiod",
            "name_space": "istio-system",
            "replicas": [
                1,
                1,
                0
            ],
            "images": "docker.io/istio/pilot:1.13.3",
            "is_complete": true,
            "key": "2"
        },
        {
            "name": "calico-kube-controllers",
            "name_space": "kube-system",
            "replicas": [
                1,
                1,
                0
            ],
            "images": "registry.cn-beijing.aliyuncs.com/kubesphereio/kube-controllers:v3.20.0",
            "is_complete": true,
            "key": "3"
        },
        {
            "name": "coredns",
            "name_space": "kube-system",
            "replicas": [
                2,
                2,
                0
            ],
            "images": "registry.cn-beijing.aliyuncs.com/kubesphereio/coredns:1.8.0",
            "is_complete": true,
            "key": "4"
        },
        {
            "name": "prodapi",
            "name_space": "myweb",
            "replicas": [
                1,
                1,
                0
            ],
            "images": "docker.io/shenyisyn/prod:v1",
            "is_complete": true,
            "key": "5"
        }
    ],
    "total": 6
}
`
	tmp := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &tmp)
	if err != nil {
		panic(err)
	}

	wsCore.ClientMap.SendAll(tmp)

}

func actionFactory(rootGroup *gin.RouterGroup) {
	ac := []handlers.Action{
		handlers.Pod{},
		handlers.Deploy{},
	}
	for _, i := range ac {
		handlers.GenRouter(rootGroup, i)
	}
}
