package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/c-bata/go-prompt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	corev1 "k8s.io/api/core/v1"
	k8slabels "k8s.io/apimachinery/pkg/labels"
)

func ShowLabels(labels map[string]string) string {
	var ret string
	for k, v := range labels {
		ret += fmt.Sprintf("%s=%s\n", k, v)
	}
	return ret
}

type CoreV1Pod []*corev1.Pod

func (c CoreV1Pod) Len() int {
	return len(c)
}

func (c CoreV1Pod) Less(i, j int) bool {
	return c[i].CreationTimestamp.Time.After(c[i].CreationTimestamp.Time)
}

func (c CoreV1Pod) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// getPodList 获取Pod的信息
func getPodList() []prompt.Suggest {
	podList, err := fact.Core().V1().Pods().Lister().Pods("default").List(k8slabels.Everything())
	if err != nil {
		log.Println(err)
	}
	sort.Sort(CoreV1Pod(podList))
	var suggests []prompt.Suggest
	for i := 0; i < len(podList); i++ {
		suggests = append(suggests, prompt.Suggest{
			Text:        podList[i].Name,
			Description: "节点:" + podList[i].Spec.NodeName + " 状态:" + string(podList[i].Status.Phase) + " IP:" + podList[i].Status.PodIP,
		})
	}
	return suggests
}

// getPodDetailByJson 打印Pod详细信息-json格式
func getPodDetailByJson(podName, path string, cmd *cobra.Command) error {
	ns, err := cmd.Flags().GetString("namespace")
	if err != nil {
		return err
	}
	if ns == "" {
		ns = "default"
	}

	pod, err := fact.Core().V1().Pods().Lister().Pods(ns).Get(podName)
	if err != nil {
		return err
	}

	// 事件
	if path == PodEventType {
		eventList, err := fact.Core().V1().Events().Lister().List(k8slabels.Everything())
		if err != nil {
			return err
		}
		var podEventList []*corev1.Event
		for _, event := range eventList {
			if event.InvolvedObject.UID == pod.UID {
				podEventList = append(podEventList, event)
			}
		}
		printEvent(podEventList)
		return nil
	}

	jsonStr, err := json.Marshal(pod)
	if err != nil {
		return err
	}

	resultData := gjson.Get(string(jsonStr), path)
	if !resultData.Exists() {
		return errors.New("无法找到对应内容, " + path)
	}

	fmt.Println(resultData.Raw)
	// 结果不是对象、数组，直接打印
	//if !resultData.IsObject() && !resultData.IsArray() {
	//	fmt.Println(resultData.Raw)
	//	return nil
	//}

	// 试试能不能进行yaml转换
	//var tempMap = map[string]interface{}{}
	//if err := yaml.Unmarshal([]byte(resultData.Raw), tempMap); err != nil {
	//	return err
	//}
	//b, err := yaml.Marshal(tempMap)
	//if err != nil {
	//	return err
	//}
	//fmt.Println(string(b))

	return nil
}

var eventHeaders = []string{"事件类型", "REASON", "所属对象", "消息"}

func printEvent(events []*corev1.Event) {
	table := tablewriter.NewWriter(os.Stdout)
	//设置头
	table.SetHeader(eventHeaders)
	for _, e := range events {
		podRow := []string{e.Type, e.Reason,
			fmt.Sprintf("%s/%s", e.InvolvedObject.Kind, e.InvolvedObject.Name), e.Message}

		table.Append(podRow)
	}
	setTable(table)
	table.Render()
}

//设置table的样式，不重要 。看看就好
func setTable(table *tablewriter.Table) {
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
}
