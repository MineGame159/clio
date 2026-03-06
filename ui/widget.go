package ui

import "github.com/gdamore/tcell/v3"

type Widget interface {
	CalcRequiredSize() (int, int)
	RequiredSize() (int, int)

	LimitSize(width, height int) (int, int)

	HandleEvent(event any)

	Draw(screen tcell.Screen, rect Rect)
}

// baseWidget

type baseWidget struct {
	requiredWidth  int
	requiredHeight int
}

func (w *baseWidget) RequiredSize() (int, int) {
	return w.requiredWidth, w.requiredHeight
}

// simpleMaxSizeWidget

type simpleMaxSizeWidget struct {
	maxWidth  int
	maxHeight int
}

func (w *simpleMaxSizeWidget) LimitSize(width, height int) (int, int) {
	if w.maxWidth > 0 {
		width = min(width, w.maxWidth)
	}
	if w.maxHeight > 0 {
		height = min(height, w.maxHeight)
	}

	return width, height
}

func (w *simpleMaxSizeWidget) SetMaxSize(width, height int) {
	w.maxWidth = width
	w.maxHeight = height
}

func (w *simpleMaxSizeWidget) SetMaxWidth(width int) {
	w.maxWidth = width
}

func (w *simpleMaxSizeWidget) SetMaxHeight(height int) {
	w.maxHeight = height
}
