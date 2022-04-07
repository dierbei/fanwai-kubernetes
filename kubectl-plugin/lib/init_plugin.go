package lib

import (
	"log"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

var (
	cfgFlags = &genericclioptions.ConfigFlags{}
	client   = InitClient()

	showLabels bool
	labels     string
	fields     string
	name       string
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

// RunCmd 运行命令
func RunCmd() error {
	cmd := &cobra.Command{
		Use:          "kubectl pods [flags]",
		Short:        "list pods",
		Example:      "kubectl pods [flags]",
		SilenceUsage: true,
	}

	// 合并K8s命令
	mergeFlags(cmd, listCmd, promptCmd)

	// 加入子命令
	cmd.AddCommand(listCmd, promptCmd)

	// 自定义参数
	cmd.Flags().BoolVar(&showLabels, "showlabels", false, "kubectl pods --showlabels 显示标签")
	cmd.Flags().StringVar(&labels, "labels", "", "kubectl pods --labels 根据标签过滤")
	cmd.Flags().StringVar(&fields, "fields", "", "kubectl pods --fields=\"status.phase=Running\"") // 参考链接: https://github.com/kubernetes/kubernetes/blob/master/pkg/apis/core/v1/conversion.go
	cmd.Flags().StringVar(&name, "name", "", "kubectl pods --name=\"^my\"")                        // github.com/tidwall/gjson

	// 执行命令
	return cmd.Execute()
}

// mergeFlags 合并K8s命令
func mergeFlags(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		cfgFlags.AddFlags(cmd.Flags())
	}
}
