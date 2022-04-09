package lib

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	k8slabels "k8s.io/apimachinery/pkg/labels"
)

var cacheCmd = &cobra.Command{
	Use:    "cache",
	Short:  "cache pods",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ns, err := cmd.Flags().GetString("namespace")
		if err != nil {
			return err
		}
		if ns == "" {
			ns = "default"
		}

		list, err := fact.Core().V1().Pods().Lister().Pods(ns).List(k8slabels.Everything())
		if err != nil {
			return err
		}
		fmt.Println("从缓存取")

		table := tablewriter.NewWriter(os.Stdout)
		var header = []string{"名称", "命名空间", "IP", "状态"}
		if showLabels {
			header = append(header, "标签")
		}
		table.SetHeader(header)

		for _, pod := range list {
			var row []string
			if name != "" {
				row, err = NameRegex(*pod, showLabels, name)

			} else {
				row = RenderData(*pod, showLabels)
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
	},
}
