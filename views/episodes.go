package views

import (
	"clio/stremio"
	"clio/ui"
	"fmt"
	"slices"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Episodes struct {
	Stack *Stack
	Ctx   *stremio.Context

	Catalog      *stremio.Catalog
	SearchResult stremio.SearchResult
	Season       uint

	input *ui.Input
	list  *ui.List[stremio.Video]
}

func (e *Episodes) Title() string {
	return fmt.Sprintf("Episodes of '%s - S%02d'", e.SearchResult.Name, e.Season)
}

func (e *Episodes) Keys() []Key {
	keys := []Key{{"Esc", "close"}}

	if e.input.Focused() {
		keys = append(keys, Key{"Enter", "search"})
	} else {
		keys = append(keys, Key{"Enter", "open"})
		keys = append(keys, Key{"Tab", "search"})
	}

	return keys
}

func (e *Episodes) Widgets() []ui.Widget {
	// List
	e.list = &ui.List[stremio.Video]{
		ItemDisplayFn: episodeWidget,
		ItemHeight:    1,
		SelectedStr:   "│ ",
		SelectedStyle: ui.Fg(color.Lime),
	}

	e.list.Focus()

	go func() {
		provider := e.Ctx.MetaProviderForId(e.SearchResult.Id)

		if provider != nil {
			if meta, err := provider.Get(e.Catalog.Type, e.SearchResult.Id); err == nil {
				e.Stack.Post(meta)
			}
		}
	}()

	// Input
	e.input = &ui.Input{
		Placeholder:      "Search episodes",
		PlaceholderStyle: ui.Fg(color.Gray),
		OnChange: func(value string) {
			e.list.Filter(ui.FilterFn(e.input.Value, episodeText))
		},
	}

	// Root
	return []ui.Widget{
		ui.LeftMargin(e.input, 2),
		&ui.HSeparator{Runes: ui.Rounded},
		e.list,
	}
}

func (e *Episodes) HandleEvent(event any) {
	switch event := event.(type) {
	case *tcell.EventKey:
		switch event.Key() {
		case tcell.KeyEnter:
			if e.input.Focused() {
				e.input.Blur()
				e.list.Focus()
			} else if item, ok := e.list.Selected(); ok {
				e.Stack.Push(&Streams{
					Stack:        e.Stack,
					Ctx:          e.Ctx,
					Catalog:      e.Catalog,
					SearchResult: e.SearchResult,
					Season:       e.Season,
					Episode:      item.Number,
					EpisodeName:  item.Name,
				})
			}

		case tcell.KeyTab:
			if e.list.Focused() {
				e.input.Focus()
				e.list.Blur()
			}

		default:
		}

	case stremio.Meta:
		e.list.SetItems(slices.Collect(event.Episodes(e.Season)))
	}
}

func episodeWidget(item stremio.Video, selected bool) ui.Widget {
	style := tcell.StyleDefault
	if selected {
		style = ui.Fg(color.Lime)
	}

	return &ui.Paragraph{Spans: []ui.Span{
		{fmt.Sprintf("%d", item.Number), ui.Fg(color.Silver)},
		{" - ", ui.Fg(color.Gray)},
		{item.Name, style},
	}}
}

func episodeText(item stremio.Video) string {
	return fmt.Sprintf("%d - %s", item.Number, item.Name)
}
