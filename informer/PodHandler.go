package informer

import "github.com/phuslu/log"

type PodHandler struct {
}

func (p PodHandler) OnAdd(obj interface{}) {
	//TODO implement me
	log.Error().Msgf("implement me")
}

func (p PodHandler) OnUpdate(oldObj, newObj interface{}) {
	//TODO implement me
	log.Error().Msgf("implement me")
}

func (p PodHandler) OnDelete(obj interface{}) {
	//TODO implement me
	log.Error().Msgf("implement me")
}
