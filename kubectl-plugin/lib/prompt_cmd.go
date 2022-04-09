package lib

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	k8slabels "k8s.io/apimachinery/pkg/labels"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

// 参考链接：https://github.com/c-bata/go-prompt
// kubectl pods prompt
var promptCmd = &cobra.Command{
	Use:          "prompt",
	Short:        "prompt pods",
	Example:      "kubectl pods prompt",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		InitSharedInformerFactory()

		p := prompt.New(
			executorCmd(cmd),
			completer,
			prompt.OptionPrefix(">>> "),
		)
		p.Run()
		return nil
	},
}

// executorCmd 装饰器模式，监视命令行输入命令
func executorCmd(cmd *cobra.Command) func(in string) {
	return func(in string) {
		in = strings.TrimSpace(in)
		blocks := strings.Split(in, " ")
		args := blocks[1:]
		switch blocks[0] {
		case "exit":
			fmt.Println("bye bye")
			os.Exit(0)
		case "list":
			if err := cacheCmd.RunE(cmd, []string{}); err != nil {
				log.Println(err)
			}
		case "get":
			if err := getPod(cmd, args); err != nil {
				log.Println(err)
			}
		}
	}
}

var suggestions = []prompt.Suggest{
	{"test", "this is test"},
	{"exit", "exit prompt"},
}

// completer 根据命令命令显示提示信息
func completer(in prompt.Document) []prompt.Suggest {
	w := in.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}
	cmd, prefix := parseCmd(in.TextBeforeCursor())
	if cmd == "get" {
		return prompt.FilterHasPrefix(getPodList(), prefix, true)
	}
	return prompt.FilterHasPrefix(suggestions, w, true)
}

// getPod 获取Pod信息，输出yaml到标准输出
func getPod(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("无效的Pod名字")
	}

	ns, err := cmd.Flags().GetString("namespace")
	if err != nil {
		return err
	}
	if ns == "" {
		ns = "default"
	}

	pod, err := fact.Core().V1().Pods().Lister().Pods(ns).Get(args[0])
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(pod)
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
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

// parseCmd 解析命令
// 例如：（get    my）==>（get my）多空格替换
func parseCmd(w string) (cmd string, suggestPrefix string) {
	w = regexp.MustCompile("\\s+").ReplaceAllString(w, " ")
	cmdSlice := strings.Split(w, " ")
	if len(cmdSlice) >= 2 {
		return cmdSlice[0], strings.Join(cmdSlice[1:], " ")
	}
	return "", ""
}

var podSuggestions = []prompt.Suggest{
	{"ngx-123", "ngx"},
	{"ngx-mygin", "mygin"},
	{"javapod-xxx", "javapod"},
}
