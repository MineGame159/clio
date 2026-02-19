package views

import (
	"clio/core"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	normalItemStyle   = lipgloss.NewStyle().Align(lipgloss.Center).PaddingLeft(3)
	selectedItemStyle = lipgloss.NewStyle().Align(lipgloss.Center).PaddingLeft(1).Foreground(core.LightGreen)
)

type SimpleItem interface {
	list.Item

	Text() string
}

type SimpleDelegate struct{}

func (s SimpleDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	render := normalItemStyle.Render

	if index == m.Index() {
		render = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprint(w, render(item.(SimpleItem).Text()))
}

func (s SimpleDelegate) Height() int {
	return 1
}

func (s SimpleDelegate) Spacing() int {
	return 0
}

func (s SimpleDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
