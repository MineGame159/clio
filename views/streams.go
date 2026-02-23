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
	list  *ui.List[Stream]
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
	s.list = &ui.List[Stream]{
		ItemDisplayFn: StreamWidget,
		ItemHeight:    2,
		SelectedStr:   "│ ",
		SelectedStyle: ui.Fg(color.Lime),
	}

	s.list.Focus()

	go func() {
		provider := s.Ctx.StreamProviderForKindId(s.Catalog.Type, s.SearchResult.Id)

		if provider != nil {
			var streams []stremio.Stream

			if s.EpisodeName != "" {
				streams, _ = provider.SearchEpisode(s.Catalog.Type, s.SearchResult.Id, s.Season, s.Episode)
			} else {
				streams, _ = provider.Search(s.Catalog.Type, s.SearchResult.Id)
			}

			streams2 := make([]Stream, 0, len(streams))

			for _, stream := range streams {
				if stream.Url != "" {
					streams2 = append(streams2, ParseStream(stream))
				}
			}

			s.Stack.Post(streams2)
		}
	}()

	// Input
	s.input = &ui.Input{
		Placeholder:      "Search streams",
		PlaceholderStyle: ui.Fg(color.Gray),
		OnChange: func(value string) {
			s.list.Filter(ui.FilterFn(value, StreamText))
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

	case []Stream:
		s.list.SetItems(event)
	}
}
