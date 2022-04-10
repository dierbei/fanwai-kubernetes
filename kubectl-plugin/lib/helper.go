package lib

import (
	"encoding/json"
	"errors"
	"fmt"

	"log"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	k8slabels "k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/yaml"
)

func ShowLabels(labels map[string]string) string {
	var ret string
	for k, v := range labels {
		ret += fmt.Sprintf("%s=%s\n", k, v)
	}
	return ret
}

// getPodList 获取Pod的信息
func getPodList() []prompt.Suggest {
	podList, err := fact.Core().V1().Pods().Lister().Pods("default").List(k8slabels.Everything())
	if err != nil {
		log.Println(err)
	}

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

	jsonStr, err := json.Marshal(pod)
	if err != nil {
		return err
	}

	fmt.Println(string(jsonStr))

	resultData := gjson.Get(string(jsonStr), path)
	if !resultData.Exists() {
		return errors.New("无法找到对应内容")
	}

	// 结果不是对象、数组，直接打印
	if !resultData.IsObject() && !resultData.IsArray() {
		fmt.Println(resultData.Raw)
		return nil
	}

	// 试试能不能进行yaml转换
	var tempMap = map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(resultData.Raw), tempMap); err != nil {
		return err
	}
	b, err := yaml.Marshal(tempMap)
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
}
