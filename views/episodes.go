package views

import (
	"clio/core"
	"clio/stremio"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Episodes struct {
	App    *core.App
	Kind   string
	Result stremio.SearchResult
	Season uint

	list list.Model
}

func (e *Episodes) Init() tea.Cmd {
	l := list.New([]list.Item{}, SimpleDelegate{}, 0, 0)

	l.DisableQuitKeybindings()
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.Styles.HelpStyle = list.DefaultStyles().HelpStyle.Padding(1, 1, 0, 1)

	e.list = l

	// Get meta
	return func() tea.Msg {
		provider := e.App.MetaProviderForId(e.Result.Id)

		if provider != nil {
			if meta, err := provider.Get(e.Kind, e.Result.Id); err == nil {
				return meta
			}
		}

		return nil
	}
}

func (e *Episodes) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return e, tea.Quit

		case tea.KeyEscape:
			if e.list.FilterState() == list.Unfiltered {
				return e, e.App.Pop()
			}

		case tea.KeyEnter:
			if e.list.FilterState() != list.Filtering {
				if episode, ok := e.list.SelectedItem().(Episode); ok {
					e.App.Push(&Streams{
						App:     e.App,
						Kind:    e.Kind,
						Result:  e.Result,
						Season:  e.Season,
						Episode: episode.Number,
					})
				}
			}
		}

	case stremio.Meta:
		var items []list.Item

		for video := range msg.Episodes(e.Season) {
			items = append(items, Episode{
				Number: video.Number,
				Name:   video.Name,
			})
		}

		cmd := e.list.SetItems(items)
		e.list.Select(0)

		// No idea, don't ask
		e.list.SetSize(e.list.Width(), e.list.Height())

		return e, cmd

	case tea.WindowSizeMsg:
		e.list.SetSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	e.list, cmd = e.list.Update(msg)

	return e, cmd
}

func (e *Episodes) View() string {
	return e.list.View()
}

// Episode

type Episode struct {
	Number uint
	Name   string
}

func (e Episode) FilterValue() string {
	return e.Name
}

func (e Episode) Text() string {
	return fmt.Sprintf("%d - %s", e.Number, e.Name)
}
