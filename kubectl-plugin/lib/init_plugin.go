package lib

import (
	"log"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

var (
	cfgFlags = &genericclioptions.ConfigFlags{}
)

// InitClient 初始化Kubernetes客户端
func InitClient() *kubernetes.Clientset {
	cfgFlags = genericclioptions.NewConfigFlags(true)
	config, err := cfgFlags.ToRawKubeConfigLoader().ClientConfig()
	if err != nil {
		log.Fatal(err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

// MergeFlags 合并命令
func MergeFlags(cmd *cobra.Command) {
	cfgFlags.AddFlags(cmd.Flags())
}
