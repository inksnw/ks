package k8sutils

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func InitClient() *kubernetes.Clientset {
	c, err := kubernetes.NewForConfig(K8sRestConfig())
	checkErr(err)
	return c
}
func InitDynamicClient() dynamic.Interface {
	ci, err := dynamic.NewForConfig(K8sRestConfig())
	checkErr(err)
	return ci
}

func K8sRestConfig() *rest.Config {
	configFile := filepath.Join(home(), ".kube/config")
	if exists(configFile) {
		config, err := clientcmd.BuildConfigFromFlags("", configFile)
		checkErr(err)
		return config
	}
	config, err := rest.InClusterConfig()
	checkErr(err)
	return config
}

func home() string {
	homePath, err := user.Current()
	checkErr(err)
	return homePath.HomeDir
}

func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
