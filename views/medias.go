package views

import (
	"clio/core"
	"clio/stremio"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Medias struct {
	App     *core.App
	Catalog *stremio.Catalog

	search textinput.Model
	list   list.Model
}

func (m *Medias) Init() tea.Cmd {
	// Search
	t := textinput.New()

	t.Placeholder = "Search catalog"
	t.Focus()

	m.search = t

	// List
	l := list.New([]list.Item{}, SimpleDelegate{}, 0, 0)

	l.DisableQuitKeybindings()
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.Styles.HelpStyle = list.DefaultStyles().HelpStyle.Padding(1, 1, 0, 1)

	m.list = l

	return nil
}

func (m *Medias) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyEscape:
			if !m.search.Focused() && m.list.FilterState() == list.Unfiltered {
				return m, m.App.Pop()
			}

		case tea.KeyTab:
			if m.search.Focused() {
				m.search.Blur()
			} else {
				m.search.Focus()
			}

		case tea.KeyEnter:
			if m.search.Focused() {
				m.search.Blur()

				return m, func() tea.Msg {
					if metas, err := m.Catalog.Search(m.search.Value()); err == nil {
						return metas
					}

					return nil
				}
			}

			if m.list.FilterState() != list.Filtering {
				if meta, ok := m.list.SelectedItem().(stremio.MetaBasic); ok {
					m.App.Push(&Streams{App: m.App, Kind: m.Catalog.Type, Meta: meta})
				}
			}
		}

	case []stremio.MetaBasic:
		items := make([]list.Item, len(msg))

		for i, meta := range msg {
			items[i] = meta
		}

		cmd := m.list.SetItems(items)
		m.list.Select(0)

		// No idea, don't ask
		m.list.SetSize(m.list.Width(), m.list.Height())

		return m, cmd

	case tea.WindowSizeMsg:
		m.search.Width = msg.Width - len(m.search.Prompt) - 1
		m.list.SetSize(msg.Width, msg.Height-2)
	}

	var cmd tea.Cmd

	if m.search.Focused() {
		m.search, cmd = m.search.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

var bottomBorderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), false, false, true, false).BorderForeground(core.White)

func (m *Medias) View() string {
	return lipgloss.JoinVertical(lipgloss.Top, bottomBorderStyle.Width(m.list.Width()).Render(m.search.View()), m.list.View())
}
