package ui

import (
	"github.com/gdamore/tcell/v3"
	"github.com/rivo/uniseg"
)

type Input struct {
	baseWidget

	Placeholder      string
	PlaceholderStyle tcell.Style

	Value      string
	ValueStyle tcell.Style

	OnChange func(value string)

	focused bool
	cursor  int
}

func (w *Input) Focus() {
	w.focused = true
}

func (w *Input) Blur() {
	w.focused = false
}

func (w *Input) Focused() bool {
	return w.focused
}

func (w *Input) cursorFromStart() int {
	return len(w.Value) - w.cursor
}

func (w *Input) CalcRequiredSize() (int, int) {
	w.maxHeight = 1

	w.requiredWidth = 10
	w.requiredHeight = 1

	return w.requiredWidth, w.requiredHeight
}

func (w *Input) HandleEvent(event any) {
	w.cursor = min(w.cursor, len(w.Value))

	if !w.focused {
		return
	}

	switch event := event.(type) {
	case *tcell.EventKey:
		switch event.Key() {
		case tcell.KeyRune:
			if w.cursor == 0 {
				w.Value += event.Str()
			} else {
				c := w.cursorFromStart()
				w.Value = w.Value[:c] + event.Str() + w.Value[c:]
			}
			if w.OnChange != nil {
				w.OnChange(w.Value)
			}

		case tcell.KeyBackspace:
			c := w.cursorFromStart()
			if c > 0 {
				w.Value = w.Value[:c-1] + w.Value[c:]

				if w.OnChange != nil {
					w.OnChange(w.Value)
				}
			}

		case tcell.KeyDelete:
			if w.cursor > 0 {
				c := w.cursorFromStart()
				w.Value = w.Value[:c] + w.Value[c+1:]
				w.cursor--

				if w.OnChange != nil {
					w.OnChange(w.Value)
				}
			}

		case tcell.KeyLeft:
			if w.cursorFromStart() > 0 {
				w.cursor++
			}

		case tcell.KeyRight:
			if w.cursor > 0 {
				w.cursor--
			}

		case tcell.KeyHome:
			w.cursor = len(w.Value)

		case tcell.KeyEnd:
			w.cursor = 0

		default:
		}
	}
}

func (w *Input) Draw(screen tcell.Screen, rect Rect) {
	w.cursor = min(w.cursor, len(w.Value))

	offset := 0
	if w.cursorFromStart() > rect.Width-2 {
		offset = w.cursorFromStart() - (rect.Width - 2)
	}

	if w.Value != "" {
		putStr(screen, rect.X, rect.Y, rect.Width, w.Value[offset:], w.ValueStyle)
	} else {
		putStr(screen, rect.X, rect.Y, rect.Width, w.Placeholder, w.PlaceholderStyle)
	}

	if rect.Height > 0 && w.focused {
		valueWidth := uniseg.StringWidth(w.Value) - offset
		screen.ShowCursor(rect.X+valueWidth-w.cursor, rect.Y)
	}
}
