package main

import (
	"fmt"
	"log"

	"githup.com/dierbei/fanwai-kubernetes/config"
	"githup.com/dierbei/fanwai-kubernetes/handler"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

func main() {
	clientSet := config.NewKubernetesConfig().InitClient()

	// 监听资源
	listWatch := cache.NewListWatchFromClient(
		clientSet.CoreV1().RESTClient(),
		"configmaps",
		"default",
		fields.Everything(),
	)

	// 创建Indexer并使用
	indexer := cache.Indexers{
		cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
		"app":                MetaAnnotationIndexFunc,
	}
	myIndexer, informer := cache.NewIndexerInformer(
		listWatch,
		&v1.ConfigMap{},
		0,
		handler.NewConfigMapHandler(),
		indexer,
	)
	stopChan := make(chan struct{})
	defer close(stopChan)
	go informer.Run(stopChan)

	// 等待缓存同步
	if !cache.WaitForCacheSync(stopChan, informer.HasSynced) {
		log.Println("sync error")
	}
	fmt.Println(myIndexer.IndexKeys("app", "cm"))

	select {}
}

// MetaAnnotationIndexFunc 根据annotations过滤
func MetaAnnotationIndexFunc(obj interface{}) ([]string, error) {
	meta, err := meta.Accessor(obj)
	if err != nil {
		return []string{""}, fmt.Errorf("object has no meta: %v", err)
	}
	if app, ok := meta.GetAnnotations()["app"]; ok {
		return []string{app}, nil
	}
	return []string{}, nil
}
