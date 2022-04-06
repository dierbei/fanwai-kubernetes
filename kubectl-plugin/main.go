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

	showlabels, err := c.Flags().GetBool("showlabels")
	if err != nil {
		return err
	}

	labels, err := c.Flags().GetString("labels")
	if err != nil {
		return err
	}

	fields, err := c.Flags().GetString("fields")
	if err != nil {
		return err
	}

	if ns == "" {
		ns = "default"
	}

	list, err := client.CoreV1().Pods(ns).List(context.Background(), v1.ListOptions{
		LabelSelector: labels,
		FieldSelector: fields,
	})
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	var header = []string{"名称", "命名空间", "IP", "状态"}
	if showlabels {
		header = append(header, "标签")
	}
	table.SetHeader(header)

	for _, pod := range list.Items {
		var row = []string{pod.Name, pod.Namespace, pod.Status.PodIP, string(pod.Status.Phase)}
		if showlabels {
			row = append(row, lib.ShowLabels(pod.Labels))
		}
		table.Append(row)
	}

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)

	table.Render() // Send output

	return nil
}
