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

// RunCmd 执行命令
func RunCmd(f func(c *cobra.Command, args []string) error) error {
	cmd := &cobra.Command{
		Use:          "kubectl pods [flags]",
		Short:        "获取Pod列表",
		Example:      "kubectl pods [flags]",
		SilenceUsage: true,
		RunE: f,
	}
	MergeFlags(cmd)
	if err := cmd.Execute(); err != nil {
		return err
	}
	return nil
}