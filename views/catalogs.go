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

	list *ui.List[catalogItem]
}

type catalogItem struct {
	addon   *stremio.Addon
	catalog *stremio.Catalog
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
	c.list = &ui.List[catalogItem]{
		ItemDisplayFn:    catalogItemWidget,
		ItemSelectableFn: catalogItemSelectable,
		ItemHeight:       1,
		SelectedStr:      "│ ",
		SelectedStyle:    ui.Fg(color.Lime),
	}

	var items []catalogItem

	for catalog := range c.Ctx.Catalogs() {
		ok := true

		for _, extra := range catalog.Extras {
			if extra.Required && extra.Name != "search" {
				ok = false
				break
			}
		}

		if ok {
			if len(items) == 0 || items[len(items)-1].addon != catalog.Addon {
				items = append(items, catalogItem{
					addon:   catalog.Addon,
					catalog: nil,
				})
			}

			items = append(items, catalogItem{
				addon:   catalog.Addon,
				catalog: catalog,
			})
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
					Catalog: item.catalog,
				})
			}

		default:
		}
	}
}

func catalogItemWidget(item catalogItem, selected bool) ui.Widget {
	// Addon
	if item.catalog == nil {
		return ui.NewParagraphStyled(item.addon.Name, ui.Fg(color.Silver).Bold(true))
	}

	// Catalog
	style := tcell.StyleDefault
	if selected {
		style = ui.Fg(color.Lime)
	}

	return &ui.Paragraph{Spans: []ui.Span{
		{"  ", style},
		{item.catalog.Kind.Name(), style},
		{" - ", ui.Fg(color.Gray)},
		{item.catalog.Name, style},
	}}
}

func catalogItemSelectable(item catalogItem) bool {
	return item.catalog != nil
}
