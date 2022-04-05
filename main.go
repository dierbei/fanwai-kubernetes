package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"githup.com/dierbei/fanwai-kubernetes/config"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func main() {
	clientSet := config.NewKubernetesConfig().InitClient()

	fact := informers.NewSharedInformerFactoryWithOptions(clientSet, 0, informers.WithNamespace("default"))
	fact.Core().V1().ConfigMaps().Informer().AddIndexers(cache.Indexers{
		"labels": CmIndexFunc,
	})
	fact.Start(wait.NeverStop)

	engine := gin.Default()

	engine.GET("/common/:gvr", func(c *gin.Context) {
		gvr, _ := schema.ParseResourceArg(c.Param("gvr"))
		informer, _ := fact.ForResource(*gvr)

		list, _ := informer.Informer().GetIndexer().
		ByIndex("labels", c.Query("labels"))
		c.JSON(200, gin.H{"data": list})
	})
	//
	//engine.GET("/:group/:version/:resource", func(ctx *gin.Context) {
	//	var set map[string]string
	//
	//	if labelsQuery, ok := ctx.GetQueryMap("labels"); ok {
	//		set = labelsQuery
	//	}
	//	var g, v, r = ctx.Param("group"), ctx.Param("version"), ctx.Param("resource")
	//	if g == "core" {
	//		g = ""
	//	}
	//
	//	gvr := schema.GroupVersionResource{Group: g, Resource: r, Version: v}
	//	informer, _ := fact.ForResource(gvr)
	//	list, _ := informer.Lister().List(labels.SelectorFromSet(set))
	//	ctx.JSON(200, gin.H{"data": list})
	//})
	//
	//engine.GET("/configmaps", func(ctx *gin.Context) {
	//	var set map[string]string
	//	if labelsQuery, ok := ctx.GetQueryMap("labels"); ok {
	//		set = labelsQuery
	//	}
	//
	//	configMaps, err := fact.Core().V1().ConfigMaps().Lister().List(labels.SelectorFromSet(set))
	//	if err != nil {
	//		ctx.JSON(http.StatusInternalServerError, gin.H{
	//			"error": err.Error(),
	//		})
	//	} else {
	//		ctx.JSON(http.StatusOK, gin.H{
	//			"result": configMaps,
	//		})
	//	}
	//})

	engine.Run(":8080")
}

func CmIndexFunc(obj interface{}) ([]string, error) {
	meta, err := meta.Accessor(obj)
	if err != nil {
		return []string{""}, fmt.Errorf("object has no meta: %v", err)
	}
	ret := []string{}
	if meta.GetLabels() != nil {
		for k, v := range meta.GetLabels() {
			//  best:true
			ret = append(ret, fmt.Sprintf("%s:%s", k, v))
		}
	}
	fmt.Println(ret)
	return ret, nil
}
