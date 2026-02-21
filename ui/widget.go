package ui

import "github.com/gdamore/tcell/v3"

type Widget interface {
	MaxSize() (int, int)

	CalcRequiredSize() (int, int)
	RequiredSize() (int, int)

	HandleEvent(event any)

	Draw(screen tcell.Screen, rect Rect)
}

// baseWidget

type baseWidget struct {
	maxWidth  int
	maxHeight int

	requiredWidth  int
	requiredHeight int
}

func (w *baseWidget) MaxSize() (int, int) {
	return w.maxWidth, w.maxHeight
}

func (w *baseWidget) RequiredSize() (int, int) {
	return w.requiredWidth, w.requiredHeight
}

func (w *baseWidget) SetMaxSize(width, height int) {
	w.maxWidth = width
	w.maxHeight = height
}

func (w *baseWidget) SetMaxWidth(width int) {
	w.maxWidth = width
}

func (w *baseWidget) SetMaxHeight(height int) {
	w.maxHeight = height
}
