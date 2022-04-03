package main

import (
	"context"
	"fmt"
	"log"

	"githup.com/dierbei/fanwai-kubernetes/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func main() {
	dynamicClient := config.NewKubernetesConfig().InitDynamicClient()
	deployGVR := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}

	deployList, err := dynamicClient.Resource(deployGVR).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
	}

	for _, deploy := range deployList.Items {
		fmt.Println(deploy.GetName())
	}
}
