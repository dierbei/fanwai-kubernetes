package main

import (
	"context"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"githup.com/dierbei/fanwai-kubernetes/kubectl-plugin/lib"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var client *kubernetes.Clientset

func main() {
	client = lib.InitClient()
	if err := lib.RunCmd(run); err != nil {
		log.Fatal(err)
	}
}

func run(c *cobra.Command, args []string) error {
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"名称", "命名空间", "IP", "状态"})
	for _, pod := range list.Items {
		table.Append([]string{pod.Name, pod.Namespace, pod.Status.PodIP, string(pod.Status.Phase)})
	}
	table.Render() // Send output

	return nil
}
