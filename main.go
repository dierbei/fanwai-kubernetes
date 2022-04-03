package main

import (
	"context"
	_ "embed"
	"log"

	"githup.com/dierbei/fanwai-kubernetes/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
)

//go:embed tpls/deploy.yaml
var deployTpl string

func main() {
	dynamicClient := config.NewKubernetesConfig().InitDynamicClient()
	deployGVR := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}

	var deployUnstructured = &unstructured.Unstructured{}
	if err := yaml.Unmarshal([]byte(deployTpl), deployUnstructured); err != nil {
		log.Println(err)
	}

	_, err := dynamicClient.Resource(deployGVR).
		Namespace(deployUnstructured.GetNamespace()).
		Create(context.Background(), deployUnstructured, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
	}
}
