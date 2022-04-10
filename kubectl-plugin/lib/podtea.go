package lib

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type podJson struct {
	title string
	path  string
}

type podModel struct {
	items   []*podJson
	index   int
	podName string
	cmd     *cobra.Command
}

func (podModel) Init() tea.Cmd {
	return nil
}

func (m podModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch t := msg.(type) {
	case tea.KeyMsg:
		switch t.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.index > 0 {
				m.index--
			}
		case "down":
			if m.index < len(m.items)-1 {
				m.index++
			}
		case "enter":
			if err := getPodDetailByJson(m.podName, m.items[m.index].path, m.cmd); err != nil {
				log.Println(err)
			}
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m podModel) View() string {
	s := "按上下键选择要查看的内容\n\n"

	for i, item := range m.items {
		cursor := " "
		if i == m.index {
			cursor = "»"
		}
		s += fmt.Sprintf("%s %s\n", cursor, item.title)
	}

	s += "\n按Q退出\n"
	return s
}

func runTea(args []string, cmd *cobra.Command) {
	if len(args) == 0 {
		log.Println("请填写有效的pod名称")
		return
	}

	var podModel = podModel{
		items:   []*podJson{},
		index:   0,
		podName: args[0],
		cmd:     cmd,
	}

	podModel.items = append(podModel.items,
		&podJson{title: "元信息", path: "metadata"},
		&podJson{title: "标签", path: "metadata.labels"},
		&podJson{title: "注解", path: "metadata.annotations"},
		&podJson{title: "容器", path: "spec"},
		&podJson{title: "全部", path: "@this"},
	)

	teaCmd := tea.NewProgram(podModel)
	if err := teaCmd.Start(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
