package informer

import "github.com/phuslu/log"

type DepHandler struct {
}

func (d DepHandler) OnAdd(obj interface{}) {
	//TODO implement me
	log.Error().Msgf("implement me")
}

func (d DepHandler) OnUpdate(oldObj, newObj interface{}) {
	//TODO implement me
	log.Error().Msgf("implement me")
}

func (d DepHandler) OnDelete(obj interface{}) {
	//TODO implement me
	log.Error().Msgf("implement me")
}
