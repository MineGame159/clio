package views

import (
	"clio/core"
	"clio/stremio"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Catalogs struct {
	App *core.App

	list list.Model
}

func (c *Catalogs) Init() tea.Cmd {
	var items []list.Item

	for catalog := range c.App.Catalogs() {
		if catalog.HasExtra("search") {
			items = append(items, catalog)
		}
	}

	l := list.New(items, SimpleDelegate{}, 0, 0)

	l.DisableQuitKeybindings()
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.Styles.HelpStyle = list.DefaultStyles().HelpStyle.Padding(1, 1, 0, 1)

	c.list = l

	return nil
}

func (c *Catalogs) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return c, tea.Quit

		case tea.KeyEscape:
			if c.list.FilterState() == list.Unfiltered {
				return c, c.App.Pop()
			}

		case tea.KeyEnter:
			c.App.Push(&Medias{
				App:     c.App,
				Catalog: c.list.SelectedItem().(*stremio.Catalog),
			})
		}

	case tea.WindowSizeMsg:
		c.list.SetSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(msg)

	return c, cmd
}

func (c *Catalogs) View() string {
	return c.list.View()
}
