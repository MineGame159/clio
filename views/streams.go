package views

import (
	"clio/core"
	"clio/stremio"
	"clio/ui"
	"context"
	"fmt"
	"net/http"

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

	providers []*stremio.StreamProvider
	providerI int

	fetchCancelFn           func()
	manuallyChangedProvider bool

	input *ui.Input
	list  *ui.List[Stream]
}

type streamCheckResult struct {
	checkUrl string
	cache    CacheStatus
}

func (s *Streams) Title() string {
	if len(s.providers) == 0 {
		return "No addon for this media kind"
	}

	return s.providers[s.providerI].Addon.Name
}

func (s *Streams) Keys() []Key {
	keys := []Key{{"Esc", "close"}}

	if s.input.Focused() {
		keys = append(keys, Key{"Enter", "search"})
	} else {
		keys = append(keys, Key{"Enter", "play"})
		keys = append(keys, Key{"Tab", "search"})

		if item, ok := s.list.Selected(); ok && item.CheckUrl != "" && item.Cache == Unknown {
			keys = append(keys, Key{"C", "check status"})
		}

		if len(s.providers) > 1 {
			keys = append(keys, Key{"B", "prev addon"})
			keys = append(keys, Key{"N", "next addon"})
		}
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

	for provider := range s.Ctx.StreamProviders() {
		if provider.SupportsKindId(s.Catalog.Kind, s.SearchResult.Id) {
			s.providers = append(s.providers, provider)
		}
	}

	if len(s.providers) > 0 {
		s.fetch()
	}

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
				url := item.Url

				if item.RedirectUrl {
					client := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
						return http.ErrUseLastResponse
					}}

					if res, err := client.Get(url); err == nil && res.StatusCode >= 300 && res.StatusCode <= 399 {
						url = res.Header.Get("Location")
					}
				}

				if s.EpisodeName != "" {
					core.OpenMpv(fmt.Sprintf("%s - S%02dE%02d - %s", s.SearchResult.Name, s.Season, s.Episode, s.EpisodeName), url)
				} else {
					core.OpenMpv(s.SearchResult.Name, url)
				}

				s.Stack.Stop()
			}

		case tcell.KeyTab:
			if s.list.Focused() {
				s.input.Focus()
				s.list.Blur()
			}

		case tcell.KeyRune:
			switch event.Str() {
			case "C", "c":
				if item := s.list.SelectedPtr(); item != nil && item.CheckUrl != "" && item.Cache == Unknown {
					item.Cache = Waiting

					go func() {
						check, err := core.GetJson[stremio.StreamCheck](item.CheckUrl)
						status := Unknown

						if err == nil {
							if check.Cached {
								status = Cached
							} else {
								status = Uncached
							}
						}

						s.Stack.Post(streamCheckResult{
							checkUrl: item.CheckUrl,
							cache:    status,
						})
					}()
				}

			case "B", "b":
				if s.list.Focused() && len(s.providers) > 1 {
					s.providerI--
					if s.providerI < 0 {
						s.providerI = len(s.providers) - 1
					}

					s.manuallyChangedProvider = true
					s.fetch()
				}

			case "N", "n":
				if s.list.Focused() && len(s.providers) > 1 {
					s.providerI = (s.providerI + 1) % len(s.providers)

					s.manuallyChangedProvider = true
					s.fetch()
				}
			}

		default:
		}

	case []Stream:
		if len(event) == 0 && !s.manuallyChangedProvider {
			s.providerI = (s.providerI + 1) % len(s.providers)

			if s.providerI == 0 {
				s.manuallyChangedProvider = false
			}

			s.fetch()
		}

		s.list.SetItems(event)

		if value := s.input.Value(); value != "" {
			s.list.Filter(ui.FilterFn(value, StreamText))
		}

	case streamCheckResult:
		if item := s.list.Modify(func(item Stream) bool {
			return item.CheckUrl == event.checkUrl
		}); item != nil {
			item.Cache = event.cache
		}
	}
}

func (s *Streams) fetch() {
	s.list.SetItems(nil)

	if s.fetchCancelFn != nil {
		s.fetchCancelFn()
		s.fetchCancelFn = nil
	}

	provider := s.providers[s.providerI]

	var ctx context.Context
	ctx, s.fetchCancelFn = context.WithCancel(context.Background())

	go func() {
		var streams []stremio.Stream

		if s.EpisodeName != "" {
			streams, _ = provider.SearchEpisode(ctx, s.Catalog.Kind, s.SearchResult.Id, s.Season, s.Episode)
		} else {
			streams, _ = provider.Search(ctx, s.Catalog.Kind, s.SearchResult.Id)
		}

		streams2 := make([]Stream, 0, len(streams))

		for _, stream := range streams {
			if stream.Url != "" {
				streams2 = append(streams2, ParseStream(stream))
			}
		}

		s.Stack.Post(streams2)
	}()
}
