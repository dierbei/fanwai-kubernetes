package config

import (
	"log"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	_kubernetesConfigPath = "resource/config"
)

type KubernetesConfig struct {
}

func NewKubernetesConfig() *KubernetesConfig {
	return &KubernetesConfig{}
}

// RestConfig 创建RestConfig
func (*KubernetesConfig) RestConfig() *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags("", _kubernetesConfigPath)
	config.Insecure = true
	if err != nil {
		log.Fatal(err)
	}
	return config
}

// InitClient 初始化客户端
func (k *KubernetesConfig) InitClient() *kubernetes.Clientset {
	client, err := kubernetes.NewForConfig(k.RestConfig())
	if err != nil {
		log.Fatal(err)
	}
	return client
}

// InitDynamicClient 初始化动态客户端-用于Operator
func (k *KubernetesConfig) InitDynamicClient() dynamic.Interface {
	client, err := dynamic.NewForConfig(k.RestConfig())
	if err != nil {
		log.Fatal(err)
	}
	return client
}
