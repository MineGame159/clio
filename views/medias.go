package views

import (
	"clio/stremio"
	"clio/ui"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Medias struct {
	Stack *Stack
	Ctx   *stremio.Context

	Catalog *stremio.Catalog

	input *ui.Input
	list  *ui.List[stremio.SearchResult]
}

func (m *Medias) Title() string {
	return m.Catalog.FullName()
}

func (m *Medias) Keys() []Key {
	keys := []Key{{"Esc", "close"}}

	if m.input.Focused() {
		keys = append(keys, Key{"Enter", "search"})
	} else {
		keys = append(keys, Key{"Enter", "open"})
		keys = append(keys, Key{"Tab", "search"})
	}

	return keys
}

func (m *Medias) Widgets() []ui.Widget {
	// List
	m.list = &ui.List[stremio.SearchResult]{
		ItemDisplayFn: ui.SimpleItemDisplayFn(searchResultText, ui.Fg(color.Lime)),
		ItemHeight:    1,
		SelectedStr:   "│ ",
		SelectedStyle: ui.Fg(color.Lime),
	}

	// Input
	m.input = &ui.Input{
		Placeholder:      "Search catalog",
		PlaceholderStyle: ui.Fg(color.Gray),
	}

	m.input.Focus()

	// Root
	return []ui.Widget{
		ui.LeftMargin(m.input, 2),
		&ui.HSeparator{Runes: ui.Rounded},
		m.list,
	}
}

func (m *Medias) HandleEvent(event any) {
	switch event := event.(type) {
	case *tcell.EventKey:
		switch event.Key() {
		case tcell.KeyEnter:
			if m.input.Focused() {
				m.list.SetItems(nil)

				go func() {
					if results, err := m.Catalog.Search(m.input.Value()); err == nil {
						m.Stack.Post(results)
					}
				}()

				m.input.Blur()
				m.list.Focus()
			} else if item, ok := m.list.Selected(); ok {
				if m.Catalog.Type == "movie" {
					m.Stack.Push(&Streams{
						Stack:        m.Stack,
						Ctx:          m.Ctx,
						Catalog:      m.Catalog,
						SearchResult: item,
						Season:       0,
						Episode:      0,
						EpisodeName:  "",
					})
				} else {
					m.Stack.Push(&Seasons{
						Stack:        m.Stack,
						Ctx:          m.Ctx,
						Catalog:      m.Catalog,
						SearchResult: item,
					})
				}
			}

		case tcell.KeyTab:
			if m.list.Focused() {
				m.input.Focus()
				m.list.Blur()
			}

		default:
		}

	case []stremio.SearchResult:
		m.list.SetItems(event)
	}
}

func searchResultText(item stremio.SearchResult) string {
	return item.Name
}
