package ui

import (
	"unicode"

	"github.com/gdamore/tcell/v3"
	"github.com/rivo/uniseg"
)

type Input struct {
	baseWidget

	Placeholder      string
	PlaceholderStyle tcell.Style

	ValueStyle tcell.Style

	OnChange func(value string)

	focused bool
	cursor  int

	runes []rune
}

func (w *Input) Value() string {
	return string(w.runes)
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

func (w *Input) CalcRequiredSize() (int, int) {
	w.maxHeight = 1

	w.requiredWidth = 10
	w.requiredHeight = 1

	return w.requiredWidth, w.requiredHeight
}

func (w *Input) HandleEvent(event any) {
	w.cursor = min(w.cursor, len(w.runes))

	if !w.focused {
		return
	}

	switch event := event.(type) {
	case *tcell.EventKey:
		switch event.Key() {
		case tcell.KeyRune:
			if w.cursor == 0 {
				w.runes = append(w.runes, []rune(event.Str())...)
			} else {
				c := w.cursorFromStart()
				w.runes = append(append(w.runes[:c], []rune(event.Str())...), w.runes[c:]...)
			}
			if w.OnChange != nil {
				w.OnChange(w.Value())
			}

		case tcell.KeyBackspace:
			c := w.cursorFromStart()

			if c > 0 {
				if event.Modifiers()&tcell.ModAlt != 0 {
					moved := w.moveLeftWord()
					w.runes = append(w.runes[:w.cursorFromStart()], w.runes[c:]...)
					w.cursor -= moved
				} else {
					w.runes = append(w.runes[:c-1], w.runes[c:]...)
				}

				if w.OnChange != nil {
					w.OnChange(w.Value())
				}
			}

		case tcell.KeyDelete:
			if w.cursor > 0 {
				c := w.cursorFromStart()

				if event.Modifiers()&tcell.ModAlt != 0 {
					w.moveRightWord()
					w.runes = append(w.runes[:c], w.runes[w.cursorFromStart():]...)
				} else {
					w.runes = append(w.runes[:c], w.runes[c+1:]...)
					w.cursor--
				}

				if w.OnChange != nil {
					w.OnChange(w.Value())
				}
			}

		case tcell.KeyLeft:
			if w.cursorFromStart() > 0 {
				if event.Modifiers()&tcell.ModAlt != 0 {
					w.moveLeftWord()
				} else {
					w.cursor++
				}
			}

		case tcell.KeyRight:
			if w.cursor > 0 {
				if event.Modifiers()&tcell.ModAlt != 0 {
					w.moveRightWord()
				} else {
					w.cursor--
				}
			}

		case tcell.KeyHome:
			w.cursor = len(w.runes)

		case tcell.KeyEnd:
			w.cursor = 0

		default:
		}
	}
}

func (w *Input) Draw(screen tcell.Screen, rect Rect) {
	w.cursor = min(w.cursor, len(w.runes))

	offset := 0
	if w.cursorFromStart() > rect.Width-2 {
		offset = w.cursorFromStart() - (rect.Width - 2)
	}

	if len(w.runes) != 0 {
		putStr(screen, rect.X, rect.Y, rect.Width, string(w.runes[offset:]), w.ValueStyle)
	} else {
		putStr(screen, rect.X, rect.Y, rect.Width, w.Placeholder, w.PlaceholderStyle)
	}

	if rect.Height > 0 && w.focused {
		valueWidth := uniseg.StringWidth(w.Value()) - offset
		screen.ShowCursor(rect.X+valueWidth-w.cursor, rect.Y)
	}
}

func (w *Input) cursorFromStart() int {
	return len(w.runes) - w.cursor
}

func (w *Input) moveLeftWord() int {
	start := w.cursor

	w.cursor++

	for w.cursorFromStart() > 0 && unicode.Is(unicode.White_Space, w.runes[w.cursorFromStart()]) {
		w.cursor++
	}

	for w.cursorFromStart() > 0 && !unicode.Is(unicode.White_Space, w.runes[w.cursorFromStart()]) && (w.cursorFromStart() == 0 || !unicode.Is(unicode.White_Space, w.runes[w.cursorFromStart()-1])) {
		w.cursor++
	}

	return w.cursor - start
}

func (w *Input) moveRightWord() int {
	start := w.cursor

	w.cursor--

	for w.cursor > 0 && unicode.Is(unicode.White_Space, w.runes[w.cursorFromStart()]) {
		w.cursor--
	}

	for w.cursor > 0 && !unicode.Is(unicode.White_Space, w.runes[w.cursorFromStart()]) && (w.cursor > 0 || !unicode.Is(unicode.White_Space, w.runes[w.cursorFromStart()+1])) {
		w.cursor--
	}

	return start - w.cursor
}
