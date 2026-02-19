package views

import (
	"clio/core"
	"clio/stremio"
	"cmp"
	"slices"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Streams struct {
	App  *core.App
	Kind string
	Meta stremio.MetaBasic

	list list.Model
}

func (s *Streams) Init() tea.Cmd {
	// List
	l := list.New([]list.Item{}, StreamDelegate{}, 0, 0)

	l.DisableQuitKeybindings()
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.Styles.HelpStyle = list.DefaultStyles().HelpStyle.Padding(1, 1, 0, 1)

	s.list = l

	// Get streams
	return func() tea.Msg {
		var provider *stremio.StreamProvider

		for streamProvider := range s.App.StreamProviders() {
			if streamProvider.SupportsId(s.Meta.Id) {
				provider = streamProvider
				break
			}
		}

		if provider != nil {
			if streams, err := provider.Search(s.Kind, s.Meta.Id); err == nil {
				return streams
			}
		}

		return nil
	}
}

func (s *Streams) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				if stream, ok := s.list.SelectedItem().(*Stream); ok {
					core.OpenMpv(s.Meta.Name, stream.Url)
					return s, tea.Quit
				}
			}
		}

	case []stremio.Stream:
		items := make([]list.Item, 0, len(msg))

		for _, stream := range msg {
			if stream.Url != "" {
				bingeGroupCount := 0

				if stream.Hints.BingeGroup == "" {
					bingeGroupCount = 1
				} else {
					for _, stream2 := range msg {
						if stream2.Hints.BingeGroup == stream.Hints.BingeGroup {
							bingeGroupCount++
						}
					}
				}

				items = append(items, &Stream{
					Name:            stream.TorrentName(),
					Resolution:      stream.Resolution(),
					Size:            stream.Size(),
					VideosInTorrent: bingeGroupCount,
					Url:             stream.Url,
				})
			}
		}

		slices.SortFunc(items, func(a, b list.Item) int {
			return cmp.Compare(b.(*Stream).Size, a.(*Stream).Size)
		})

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

func (s *Streams) View() string {
	return s.list.View()
}
