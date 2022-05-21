package informer

import (
	wsCore "github.com/ks/rest/websocket"
)

type Common struct {
}

func (d Common) OnAdd(obj interface{}) {
	wsCore.ClientMap.SendAll("deployment add")
}

func (d Common) OnUpdate(oldObj, newObj interface{}) {
	wsCore.ClientMap.SendAll("deployment update")
}

func (d Common) OnDelete(obj interface{}) {
	wsCore.ClientMap.SendAll("deployment delete")
}
