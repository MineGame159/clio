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

	requestedMetaId string

	input *ui.Input
	list  *ui.List[stremio.SearchResult]

	metaParagraph *ui.Paragraph
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
	// Input
	m.input = &ui.Input{
		Placeholder:      "Search catalog",
		PlaceholderStyle: ui.Fg(color.Gray),
	}

	m.input.Focus()

	// List
	m.list = &ui.List[stremio.SearchResult]{
		ItemDisplayFn: ui.SimpleItemDisplayFn(searchResultText, ui.Fg(color.Lime)),
		ItemHeight:    1,
		SelectedStr:   "│ ",
		SelectedStyle: ui.Fg(color.Lime),
	}

	// Meta
	m.metaParagraph = &ui.Paragraph{}
	m.metaParagraph.SetMaxWidth(100)

	// Root
	return []ui.Widget{
		ui.LeftMargin(m.input, 2),
		&ui.HSeparator{Runes: ui.Rounded},
		&ui.Container{
			Direction:          ui.Horizontal,
			PrimaryAlignment:   ui.Stretch,
			SecondaryAlignment: ui.Stretch,
			Children: []ui.Widget{
				m.list,
				&ui.VSeparator{Runes: ui.Rounded},
				m.metaParagraph,
			},
		},
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

	case *tcell.EventResize:
		width, _ := event.Size()
		m.metaParagraph.SetMaxWidth(width / 3)

	case []stremio.SearchResult:
		m.list.SetItems(event)

	case stremio.Meta:
		var spans []ui.Span

		released := addMetaInfo(&spans, "Released", []string{event.ReleaseInfo})
		genres := addMetaInfo(&spans, "Genres", event.Genres)
		if released || genres {
			addMetaNewLine(&spans)
		}

		rating := addMetaInfo(&spans, "Rating", []string{event.Rating})
		awards := addMetaInfo(&spans, "Awards", []string{event.Awards})
		if rating || awards {
			addMetaNewLine(&spans)
		}

		status := addMetaInfo(&spans, "Status", []string{event.Status})
		runtime := addMetaInfo(&spans, "Runtime", []string{event.Runtime})
		if status || runtime {
			addMetaNewLine(&spans)
		}

		if addMetaInfo(&spans, "Cast", event.Cast) {
			addMetaNewLine(&spans)
		}

		addMetaInfo(&spans, "Description", []string{event.Description})

		m.metaParagraph.Spans = spans
	}

	if item, ok := m.list.Selected(); ok && item.Id != m.requestedMetaId {
		m.requestedMetaId = item.Id

		go func() {
			provider := m.Ctx.MetaProviderForId(item.Id)

			if provider != nil {
				if meta, err := provider.Get(m.Catalog.Type, item.Id); err == nil {
					m.Stack.Post(meta)
				}
			}
		}()
	}
}

func addMetaInfo(spans *[]ui.Span, name string, values []string) bool {
	if len(values) == 0 || (len(values) == 1 && values[0] == "") {
		return false
	}

	if len(*spans) > 0 {
		addMetaNewLine(spans)
	}

	*spans = append(*spans, ui.Span{Text: name, Style: tcell.StyleDefault})
	*spans = append(*spans, ui.Span{Text: ": ", Style: ui.Fg(color.Silver)})

	for i, value := range values {
		if i > 0 {
			*spans = append(*spans, ui.Span{Text: ", ", Style: ui.Fg(color.Silver)})
		}

		*spans = append(*spans, ui.Span{Text: value, Style: ui.Fg(color.Gray)})
	}

	return true
}

func addMetaNewLine(spans *[]ui.Span) {
	*spans = append(*spans, ui.Span{Text: "\n", Style: tcell.StyleDefault})
}

func searchResultText(item stremio.SearchResult) string {
	return item.Name
}
