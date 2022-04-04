package main

import (
	"githup.com/dierbei/fanwai-kubernetes/config"
	"githup.com/dierbei/fanwai-kubernetes/handler"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
)

func main() {
	clientSet := config.NewKubernetesConfig().InitClient()

	listWatch := cache.NewListWatchFromClient(
		clientSet.CoreV1().RESTClient(),
		"configmaps",
		"default",
		fields.Everything(),
	)

	//_, infomer := cache.NewInformer(listWatch, &v1.ConfigMap{}, 0, handler.NewConfigMapHandler())
	//infomer.Run(wait.NeverStop)

	sharedInformer := cache.NewSharedInformer(listWatch, &v1.ConfigMap{}, 0)
	sharedInformer.AddEventHandler(handler.NewConfigMapHandler())
	sharedInformer.AddEventHandler(handler.NewConfigMapV2Handler())
	sharedInformer.Run(wait.NeverStop)
	select {}
}
