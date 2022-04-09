package lib

import (
	"k8s.io/client-go/informers"
)

var fact informers.SharedInformerFactory

func InitSharedInformerFactory() {
	fact = informers.NewSharedInformerFactoryWithOptions(client, 0)
	fact.Core().V1().Pods().Informer().AddEventHandler(&PodHandler{})
	ch := make(chan struct{})
	fact.Start(ch)
	fact.WaitForCacheSync(ch)
}

type PodHandler struct {
}

func (h *PodHandler) OnAdd(obj interface{}) {

}

func (h *PodHandler) OnUpdate(oldObj, newObj interface{}) {

}

func (h *PodHandler) OnDelete(obj interface{}) {

}
