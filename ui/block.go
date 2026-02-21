package ui

import "github.com/gdamore/tcell/v3"

type Block struct {
	Container

	Runes *BorderRunes
	Style tcell.Style

	Title  []Span
	Footer []Span
}

func (w *Block) CalcRequiredSize() (int, int) {
	w.padding = 1
	return w.Container.CalcRequiredSize()
}

func (w *Block) Draw(screen tcell.Screen, rect Rect) {
	// Corners
	screen.Put(rect.X, rect.Y, w.Runes.CornerTopLeft, w.Style)
	screen.Put(rect.X+rect.Width-1, rect.Y, w.Runes.CornerTopRight, w.Style)
	screen.Put(rect.X+rect.Width-1, rect.Y+rect.Height-1, w.Runes.CornerBottomRight, w.Style)
	screen.Put(rect.X, rect.Y+rect.Height-1, w.Runes.CornerBottomLeft, w.Style)

	// Horizontal lines
	for i := range rect.Width - 2 {
		screen.Put(rect.X+1+i, rect.Y, w.Runes.Horizontal, w.Style)
		screen.Put(rect.X+1+i, rect.Y+rect.Height-1, w.Runes.Horizontal, w.Style)
	}

	// Vertical lines
	for i := range rect.Height - 2 {
		screen.Put(rect.X, rect.Y+1+i, w.Runes.Vertical, w.Style)
		screen.Put(rect.X+rect.Width-1, rect.Y+1+i, w.Runes.Vertical, w.Style)
	}

	// Title
	if len(w.Title) > 0 {
		p := Paragraph{Spans: w.Title}
		p.CalcRequiredSize()

		p.Draw(screen, Rect{
			X:      rect.X + 2,
			Y:      rect.Y,
			Width:  max(rect.Width-4, 0),
			Height: 1,
		})
	}

	// Footer
	if len(w.Footer) > 0 {
		p := Paragraph{Spans: w.Footer}
		p.CalcRequiredSize()

		p.Draw(screen, Rect{
			X:      rect.X + 2,
			Y:      rect.Y + rect.Height - 1,
			Width:  max(rect.Width-4, 0),
			Height: 1,
		})
	}

	// Children
	onChildDraw := func(child Widget, rect Rect) {
		if _, ok := child.(*HSeparator); ok {
			// Vertical tees
			screen.Put(rect.X-1, rect.Y, w.Runes.TeeLeft, w.Style)
			screen.Put(rect.X+rect.Width, rect.Y, w.Runes.TeeRight, w.Style)
		} else if _, ok := child.(*VSeparator); ok {
			// Horizontal tees
			screen.Put(rect.X, rect.Y-1, w.Runes.TeeTop, w.Style)
			screen.Put(rect.X, rect.Y+rect.Height, w.Runes.TeeBottom, w.Style)
		}
	}

	if w.Padding > 0 {
		onChildDraw = nil
	}

	w.Container.draw(screen, rect, onChildDraw)
}
