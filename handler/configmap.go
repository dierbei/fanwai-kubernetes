package handler

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

type ConfigMapHandler struct {
}

func NewConfigMapHandler() *ConfigMapHandler {
	return &ConfigMapHandler{}
}

func (*ConfigMapHandler) OnAdd(obj interface{}) {
	fmt.Println("add: ", obj.(*v1.ConfigMap).Name)
}

func (*ConfigMapHandler) OnUpdate(oldObj, newObj interface{}) {
}

func (*ConfigMapHandler) OnDelete(obj interface{}) {
	fmt.Println("delete: ", obj.(*v1.ConfigMap).Name)
}
