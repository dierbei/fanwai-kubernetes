package main

import (
	"context"
	"fmt"
	"log"

	"githup.com/dierbei/fanwai-kubernetes/config"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/json"
)

func main() {
	dynamicClient := config.NewKubernetesConfig().InitDynamicClient()
	deployGVR := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}

	deployUnstructuredList, err := dynamicClient.Resource(deployGVR).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
	}

	b, err := deployUnstructuredList.MarshalJSON()
	if err != nil {
		log.Println(err)
	}

	var deployList = &v1.DeploymentList{}
	if err := json.Unmarshal(b, deployList); err != nil {
		log.Println(err)
	}

	for _, deploy := range deployList.Items {
		fmt.Println(deploy.GetName())
	}
}
