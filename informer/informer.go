package informer

import (
	"github.com/ks/k8sutils"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
)

func InitInformer() informers.SharedInformerFactory {
	fact := informers.NewSharedInformerFactory(k8sutils.Client, 0)

	depInformer := fact.Apps().V1().Deployments()
	depInformer.Informer().AddEventHandler(DepHandler{})

	podInformer := fact.Core().V1().Pods() //监听pod
	podInformer.Informer().AddEventHandler(PodHandler{})

	fact.Start(wait.NeverStop)

	return fact
}
