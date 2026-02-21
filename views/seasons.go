package views

import (
	"clio/core"
	"clio/stremio"
	"clio/ui"
	"fmt"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Seasons struct {
	Stack *Stack
	Ctx   *stremio.Context

	Catalog      *stremio.Catalog
	SearchResult stremio.SearchResult

	input *ui.Input
	list  *ui.List[Season]
}

type Season struct {
	Number   uint
	Episodes uint
}

func (s *Seasons) Title() string {
	return fmt.Sprintf("Seasons of '%s'", s.SearchResult.Name)
}

func (s *Seasons) Keys() []Key {
	keys := []Key{{"Esc", "close"}}

	if s.input.Focused() {
		keys = append(keys, Key{"Enter", "search"})
	} else {
		keys = append(keys, Key{"Enter", "open"})
		keys = append(keys, Key{"Tab", "search"})
	}

	return keys
}

func (s *Seasons) Widgets() []ui.Widget {
	// List
	s.list = &ui.List[Season]{
		ItemDisplayFn: seasonWidget,
		ItemHeight:    1,
		SelectedStr:   "│ ",
		SelectedStyle: ui.Fg(color.Lime),
	}

	s.list.Focus()

	go func() {
		provider := s.Ctx.MetaProviderForId(s.SearchResult.Id)

		if provider != nil {
			if meta, err := provider.Get(s.Catalog.Type, s.SearchResult.Id); err == nil {
				s.Stack.Post(meta)
			}
		}
	}()

	// Input
	s.input = &ui.Input{
		Placeholder:      "Search seasons",
		PlaceholderStyle: ui.Fg(color.Gray),
		OnChange: func(value string) {
			s.list.Filter(ui.FilterFn(value, seasonText))
		},
	}

	// Root
	return []ui.Widget{
		ui.LeftMargin(s.input, 2),
		&ui.HSeparator{Runes: ui.Rounded},
		s.list,
	}
}

func (s *Seasons) HandleEvent(event any) {
	switch event := event.(type) {
	case *tcell.EventKey:
		switch event.Key() {
		case tcell.KeyEnter:
			if s.input.Focused() {
				s.input.Blur()
				s.list.Focus()
			} else if item, ok := s.list.Selected(); ok {
				s.Stack.Push(&Episodes{
					Stack:        s.Stack,
					Ctx:          s.Ctx,
					Catalog:      s.Catalog,
					SearchResult: s.SearchResult,
					Season:       item.Number,
				})
			}

		case tcell.KeyTab:
			if s.list.Focused() {
				s.input.Focus()
				s.list.Blur()
			}

		default:
		}

	case stremio.Meta:
		var seasons []Season

		for _, number := range event.Seasons() {
			seasons = append(seasons, Season{
				Number:   number,
				Episodes: core.Count(event.Episodes(number)),
			})
		}

		s.list.SetItems(seasons)

		if len(seasons) > 0 && seasons[0].Number == 0 {
			s.list.Select(1)
		}
	}
}

func seasonWidget(item Season, selected bool) ui.Widget {
	style := tcell.StyleDefault
	if selected {
		style = ui.Fg(color.Lime)
	}

	return &ui.Paragraph{Spans: []ui.Span{
		{"Season ", ui.Fg(color.Silver)},
		{fmt.Sprintf("%d", item.Number), style},
		{fmt.Sprintf(" [%d episodes]", item.Episodes), ui.Fg(color.Gray)},
	}}
}

func seasonText(item Season) string {
	return fmt.Sprintf("Season %d", item.Number)
}
