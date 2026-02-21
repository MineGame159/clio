package views

import (
	"clio/stremio"
	"clio/ui"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Catalogs struct {
	Stack *Stack
	Ctx   *stremio.Context

	list *ui.List[*stremio.Catalog]
}

func (c *Catalogs) Title() string {
	return "Catalogs"
}

func (c *Catalogs) Keys() []Key {
	return []Key{
		{"Esc", "close"},
		{"Enter", "open"},
	}
}

func (c *Catalogs) Widgets() []ui.Widget {
	c.list = &ui.List[*stremio.Catalog]{
		ItemDisplayFn: ui.SimpleItemDisplayFn(catalogText, ui.Fg(color.Lime)),
		ItemHeight:    1,
		SelectedStr:   "│ ",
		SelectedStyle: ui.Fg(color.Lime),
	}

	var items []*stremio.Catalog

	for catalog := range c.Ctx.Catalogs() {
		if catalog.HasExtra("search") {
			items = append(items, catalog)
		}
	}

	c.list.SetItems(items)
	c.list.Focus()

	return []ui.Widget{c.list}
}

func (c *Catalogs) HandleEvent(event any) {
	switch event := event.(type) {
	case *tcell.EventKey:
		switch event.Key() {
		case tcell.KeyEnter:
			if item, ok := c.list.Selected(); ok {
				c.Stack.Push(&Medias{
					Stack:   c.Stack,
					Ctx:     c.Ctx,
					Catalog: item,
				})
			}

		default:
		}
	}
}

func catalogText(item *stremio.Catalog) string {
	return item.FullName()
}
