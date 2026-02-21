package views

import (
	"clio/core"
	"clio/stremio"
	"clio/ui"
	"fmt"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Streams struct {
	Stack *Stack
	Ctx   *stremio.Context

	Catalog      *stremio.Catalog
	SearchResult stremio.SearchResult
	Season       uint
	Episode      uint
	EpisodeName  string

	input *ui.Input
	list  *ui.List[stremio.Stream]
}

func (s *Streams) Title() string {
	if s.EpisodeName != "" {
		return fmt.Sprintf("Streams for '%s - S%02dE%02d - %s'", s.SearchResult.Name, s.Season, s.Episode, s.EpisodeName)
	}

	return fmt.Sprintf("Streams for '%s'", s.SearchResult.Name)
}

func (s *Streams) Keys() []Key {
	keys := []Key{{"Esc", "close"}}

	if s.input.Focused() {
		keys = append(keys, Key{"Enter", "search"})
	} else {
		keys = append(keys, Key{"Enter", "play"})
		keys = append(keys, Key{"Tab", "search"})
	}

	return keys
}

func (s *Streams) Widgets() []ui.Widget {
	// List
	s.list = &ui.List[stremio.Stream]{
		ItemDisplayFn: streamWidget,
		ItemHeight:    2,
		SelectedStr:   "│ ",
		SelectedStyle: ui.Fg(color.Lime),
	}

	s.list.Focus()

	go func() {
		provider := s.Ctx.StreamProviderForId(s.SearchResult.Id)

		if provider != nil {
			if s.EpisodeName != "" {
				if streams, err := provider.SearchEpisode(s.Catalog.Type, s.SearchResult.Id, s.Season, s.Episode); err == nil {
					s.Stack.Post(streams)
				}
			} else {
				if streams, err := provider.Search(s.Catalog.Type, s.SearchResult.Id); err == nil {
					s.Stack.Post(streams)
				}
			}
		}
	}()

	// Input
	s.input = &ui.Input{
		Placeholder:      "Search streams",
		PlaceholderStyle: ui.Fg(color.Gray),
		OnChange: func(value string) {
			s.list.Filter(ui.FilterFn(value, streamText))
		},
	}

	// Root
	return []ui.Widget{
		ui.LeftMargin(s.input, 2),
		&ui.HSeparator{Runes: ui.Rounded},
		s.list,
	}
}

func (s *Streams) HandleEvent(event any) {
	switch event := event.(type) {
	case *tcell.EventKey:
		switch event.Key() {
		case tcell.KeyEnter:
			if s.input.Focused() {
				s.input.Blur()
				s.list.Focus()
			} else if item, ok := s.list.Selected(); ok {
				if s.EpisodeName != "" {
					core.OpenMpv(fmt.Sprintf("%s - S%02dE%02d - %s", s.SearchResult.Name, s.Season, s.Episode, s.EpisodeName), item.Url)
				} else {
					core.OpenMpv(s.SearchResult.Name, item.Url)
				}

				s.Stack.Stop()
			}

		case tcell.KeyTab:
			if s.list.Focused() {
				s.input.Focus()
				s.list.Blur()
			}

		default:
		}

	case []stremio.Stream:
		s.list.SetItems(event)
	}
}

func streamWidget(item stremio.Stream, selected bool) ui.Widget {
	style := tcell.StyleDefault
	if selected {
		style = ui.Fg(color.Lime)
	}

	spans := []ui.Span{{item.TorrentName() + "\n", style}}

	resolution := item.Resolution()
	if resolution != "" {
		spans = append(spans, ui.Span{Text: resolution, Style: ui.Fg(color.Gray)})
	}

	if len(spans) > 1 {
		spans = append(spans, ui.Span{Text: ", ", Style: ui.Fg(color.Silver)})
	}
	spans = append(spans, ui.Span{Text: item.Size().String(), Style: ui.Fg(color.Gray)})

	return &ui.Paragraph{Spans: spans}
}

func streamText(item stremio.Stream) string {
	return item.TorrentName()
}
