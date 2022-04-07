package lib

import (
	"context"
	"encoding/json"
	"os"
	"regexp"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// kubectl pods list
var listCmd = &cobra.Command{
	Use:          "list",
	Short:        "list pods",
	Example:      "kubectl pods list",
	SilenceUsage: true,
	RunE:         run,
}

func run(c *cobra.Command, args []string) error {
	ns, err := c.Flags().GetString("namespace")
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
	if showLabels {
		header = append(header, "标签")
	}
	table.SetHeader(header)

	for _, pod := range list.Items {
		var row []string
		if name != "" {
			row, err = NameRegex(pod, showLabels, name)

		} else {
			row = RenderData(pod, showLabels)
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

// NameRegex 根据正则组装Pod数据
func NameRegex(pod corev1.Pod, showlabels bool, name string) ([]string, error) {
	var row []string

	b, err := json.Marshal(pod)
	if err != nil {
		return nil, err
	}
	ret := gjson.Get(string(b), "metadata.name")

	for _, r := range ret.Array() {
		if m, err := regexp.MatchString(name, r.String()); err == nil && m {
			row = RenderData(pod, showlabels)
		}
	}

	return row, nil
}

// RenderData 组装Pod数据
func RenderData(pod corev1.Pod, showlabels bool) []string {
	var row = []string{pod.Name, pod.Namespace, pod.Status.PodIP, string(pod.Status.Phase)}
	if showlabels {
		row = append(row, ShowLabels(pod.Labels))
	}
	return row
}
