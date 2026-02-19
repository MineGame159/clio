package views

import (
	"clio/core"
	"clio/stremio"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Seasons struct {
	App    *core.App
	Kind   string
	Result stremio.SearchResult

	list list.Model
}

func (s *Seasons) Init() tea.Cmd {
	l := list.New([]list.Item{}, SimpleDelegate{}, 0, 0)

	l.DisableQuitKeybindings()
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.Styles.HelpStyle = list.DefaultStyles().HelpStyle.Padding(1, 1, 0, 1)

	s.list = l

	// Get meta
	return func() tea.Msg {
		provider := s.App.MetaProviderForId(s.Result.Id)

		if provider != nil {
			if meta, err := provider.Get(s.Kind, s.Result.Id); err == nil {
				return meta
			}
		}

		return nil
	}
}

func (s *Seasons) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return s, tea.Quit

		case tea.KeyEscape:
			if s.list.FilterState() == list.Unfiltered {
				return s, s.App.Pop()
			}

		case tea.KeyEnter:
			if s.list.FilterState() != list.Filtering {
				if season, ok := s.list.SelectedItem().(Season); ok {
					s.App.Push(&Episodes{
						App:    s.App,
						Kind:   s.Kind,
						Result: s.Result,
						Season: uint(season),
					})
				}
			}
		}

	case stremio.Meta:
		seasons := msg.Seasons()
		items := make([]list.Item, seasons)

		for i := range seasons {
			items[i] = Season(i)
		}

		cmd := s.list.SetItems(items)
		s.list.Select(0)

		// No idea, don't ask
		s.list.SetSize(s.list.Width(), s.list.Height())

		return s, cmd

	case tea.WindowSizeMsg:
		s.list.SetSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)

	return s, cmd
}

func (s *Seasons) View() string {
	return s.list.View()
}

// Season

type Season uint

func (s Season) FilterValue() string {
	return s.Text()
}

func (s Season) Text() string {
	return fmt.Sprintf("Season %d", s)
}
