package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"githup.com/dierbei/fanwai-kubernetes/kubectl-plugin/lib"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

func main() {
	client := lib.InitClient()
	cmd := &cobra.Command{
		Use:          "kubectl pods [flags]",
		Short:        "获取Pod列表",
		Example:      "kubectl pods [flags]",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			ns, err := c.Flags().GetString("namespace")
			if err != nil {
				return err
			}

			if ns == "" {
				ns = "default"
			}

			list, err := client.CoreV1().Pods(ns).List(context.Background(), v1.ListOptions{})
			if err != nil {
				return err
			}
			for _, pod := range list.Items {
				fmt.Println(pod.Name)
			}

			return nil
		},
	}

	lib.MergeFlags(cmd)
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}