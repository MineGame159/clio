package views

import (
	"clio/stremio"
	"clio/ui"
	"fmt"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Files struct {
	Stack *Stack
	Ctx   *stremio.Context

	Name  string
	Files []stremio.File
}

func (f *Files) Title() string {
	return fmt.Sprintf("Files of '%s'", f.Name)
}

func (f *Files) Keys() []Key {
	return []Key{{"Esc", "close"}}
}

func (f *Files) Widgets() []ui.Widget {
	list := &ui.List[stremio.File]{
		ItemDisplayFn: fileWidget,
		ItemHeight:    2,
		SelectedStr:   "│ ",
		SelectedStyle: ui.Fg(color.Lime),
	}

	list.Focus()
	list.SetItems(f.Files)

	return []ui.Widget{list}
}

func (f *Files) HandleEvent(_ any) {}

func fileWidget(file stremio.File, selected bool) ui.Widget {
	style := tcell.StyleDefault
	if selected {
		style = ui.Fg(color.Lime)
	}

	return &ui.Paragraph{Spans: []ui.Span{
		{file.Path, style},
		{"\n", tcell.StyleDefault},
		{file.Size.String(), ui.Fg(color.Gray)},
	}}
}
