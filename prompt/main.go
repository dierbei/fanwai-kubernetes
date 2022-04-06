package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
)

func main() {
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
	)
	p.Run()
}

// executor 命令执行器
func executor(in string) {
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")
	switch blocks[0] {
	case "exit":
		fmt.Println("Bye!")
		os.Exit(0)
	}
}

var suggestions = []prompt.Suggest{
	{"test", "this is test"},
	{"exit", "Exit prompt"},
}

// completer 输入命令的显示提示信息
func completer(in prompt.Document) []prompt.Suggest {
	w := in.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}
	return prompt.FilterHasPrefix(suggestions, w, true)
}
