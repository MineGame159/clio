package ui

import "github.com/gdamore/tcell/v3"

// Horizontal

type HSeparator struct {
	baseWidget

	Runes *BorderRunes
	Style tcell.Style
}

func (w *HSeparator) CalcRequiredSize() (int, int) {
	w.maxHeight = 1

	w.requiredWidth = 1
	w.requiredHeight = 1

	return w.requiredWidth, w.requiredHeight
}

func (w *HSeparator) HandleEvent(_ any) {
}

func (w *HSeparator) Draw(screen tcell.Screen, rect Rect) {
	for x := range rect.Width {
		screen.Put(rect.X+x, rect.Y, w.Runes.Horizontal, w.Style)
	}
}

// Vertical

type VSeparator struct {
	baseWidget

	Runes *BorderRunes
	Style tcell.Style
}

func (w *VSeparator) CalcRequiredSize() (int, int) {
	w.maxWidth = 1

	w.requiredWidth = 1
	w.requiredHeight = 1

	return w.requiredWidth, w.requiredHeight
}

func (w *VSeparator) HandleEvent(_ any) {
}

func (w *VSeparator) Draw(screen tcell.Screen, rect Rect) {
	for y := range rect.Height {
		screen.Put(rect.X, rect.Y+y, w.Runes.Vertical, w.Style)
	}
}
