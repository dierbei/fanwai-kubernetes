package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	choices  []string           // 列表
	cursor   int                // 光标位置
	selected map[int]struct{}   // 选中集合
}

func initialModel() model {
	return model{
		choices:  []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		// 实际按下的按键
		switch msg.String() {

		// 退出
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys 向上移动光标
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys 向下移动光标
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// 更新选中集合
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// 返回选中集合为了及时更新
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// 头部标题
	s := "What should we buy at the market?\n\n"

	// 遍历列表
	for i, choice := range m.choices {

		// 判断光标是否在这个选项上
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// 判断集合中是否有该选项
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// 渲染行
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// 这是底部提示
	s += "\nPress q to quit.\n"

	// 发送数据给UI
	return s
}