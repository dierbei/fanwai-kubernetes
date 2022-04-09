package lib

import (
	"fmt"
	"log"
	"os"
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
		switch blocks[0] {
		case "exit":
			fmt.Println("bye bye")
			os.Exit(0)
		case "list":
			InitSharedInformerFactory()
			if err := cacheCmd.RunE(cmd, []string{}); err != nil {
				log.Fatal(err)
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
	return prompt.FilterHasPrefix(suggestions, w, true)
}
