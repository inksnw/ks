package informer

import (
	wsCore "awesomeProject/api/websocket"
	"github.com/phuslu/log"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/util/json"
)

type CommonEvent struct {
	Resource string
}

func (d CommonEvent) OnAdd(obj interface{}) {
	//fmt.Println(obj)
	realObj := obj.(runtime.Object)
	decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj1 := &unstructured.Unstructured{}
	gvk := realObj.GetObjectKind().GroupVersionKind()
	b, err := json.Marshal(obj)
	if err != nil {
		log.Error().Msgf("%s", err)
	}
	_, _, err = decoder.Decode(b, &gvk, obj1)
	if err != nil {
		//log.Error().Msgf("%s", err)
	}
	//log.Debug().Msgf("[%s] [%s] [%s] on add", obj1.GetNamespace(), d.Resource, obj1.GetName())
	wsCore.ClientMap.SendAll("deployment add")

}

func (d CommonEvent) OnUpdate(oldObj, newObj interface{}) {
	//log.Debug().Msgf("on update")
	wsCore.ClientMap.SendAll("deployment update")
}

func (d CommonEvent) OnDelete(obj interface{}) {
	log.Debug().Msgf("on delete")
	wsCore.ClientMap.SendAll("deployment delete")
}
