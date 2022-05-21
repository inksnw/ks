package informer

import (
	wsCore "github.com/ks/rest/websocket"
)

type DepHandler struct {
}

func (d DepHandler) OnAdd(obj interface{}) {
	wsCore.ClientMap.SendAll("deployment add")
}

func (d DepHandler) OnUpdate(oldObj, newObj interface{}) {
	wsCore.ClientMap.SendAll("deployment update")
}

func (d DepHandler) OnDelete(obj interface{}) {
	wsCore.ClientMap.SendAll("deployment delete")
}
