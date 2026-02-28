package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/rivo/uniseg"
)

type List[T any] struct {
	baseWidget

	ItemDisplayFn func(item T, selected bool) Widget
	ItemHeight    int

	SelectedStr   string
	SelectedStyle tcell.Style

	items    []T
	filtered []filteredItem

	focused  bool
	selected int

	itemsPerPage int
	pages        int
}

type filteredItem struct {
	index  int
	widget Widget
}

func SimpleItemDisplayFn[T any](itemTextFn func(item T) string, selectedStyle tcell.Style) func(item T, selected bool) Widget {
	return func(item T, selected bool) Widget {
		style := tcell.StyleDefault
		if selected {
			style = selectedStyle
		}

		return NewParagraphStyled(itemTextFn(item), style)
	}
}

func (w *List[T]) SetItems(items []T) {
	w.items = items
	w.filtered = make([]filteredItem, len(items))

	for i := range len(items) {
		w.filtered[i].index = i
	}

	w.selected = 0
	w.itemsPerPage = 0
	w.pages = 0
}

func (w *List[T]) Modify(predicate func(item T) bool) *T {
	for i := range w.items {
		if predicate(w.items[i]) {
			for j := range w.filtered {
				if w.filtered[j].index == i {
					w.filtered[j].widget = nil
					break
				}
			}

			return &w.items[i]
		}
	}

	return nil
}

func (w *List[T]) Filter(predicate func(item T) bool) {
	w.filtered = nil

	for i, item := range w.items {
		if predicate(item) {
			w.filtered = append(w.filtered, filteredItem{i, nil})
		}
	}

	w.selected = 0
	w.itemsPerPage = 0
	w.pages = 0
}

func (w *List[T]) Select(index int) {
	if index < len(w.filtered) {
		w.selected = index
	}
}

func (w *List[T]) Selected() (T, bool) {
	if w.selected >= len(w.filtered) {
		var empty T
		return empty, false
	}

	return w.items[w.filtered[w.selected].index], true
}

func (w *List[T]) SelectedPtr() *T {
	if w.selected >= len(w.filtered) {
		return nil
	}

	filtered := &w.filtered[w.selected]
	filtered.widget = nil

	return &w.items[filtered.index]
}

func (w *List[T]) Focus() {
	w.focused = true

	if len(w.filtered) > 0 {
		w.filtered[w.selected].widget = nil
	}
}

func (w *List[T]) Blur() {
	w.focused = false

	if len(w.filtered) > 0 {
		w.filtered[w.selected].widget = nil
	}
}

func (w *List[T]) Focused() bool {
	return w.focused
}

func (w *List[T]) CalcRequiredSize() (int, int) {
	w.requiredWidth = 5
	w.requiredHeight = w.ItemHeight * 2

	return w.requiredWidth, w.requiredHeight
}

func (w *List[T]) HandleEvent(event any) {
	if !w.focused {
		return
	}

	switch event := event.(type) {
	case *tcell.EventKey:
		switch event.Key() {
		case tcell.KeyUp:
			if w.selected > 0 {
				w.filtered[w.selected].widget = nil
				w.filtered[w.selected-1].widget = nil

				w.selected--
			}

		case tcell.KeyDown:
			if w.selected < len(w.filtered)-1 {
				w.filtered[w.selected].widget = nil
				w.filtered[w.selected+1].widget = nil

				w.selected++
			}

		case tcell.KeyLeft:
			if w.itemsPerPage > 0 && w.selected >= w.itemsPerPage {
				w.filtered[w.selected].widget = nil
				w.filtered[w.selected-w.itemsPerPage].widget = nil

				w.selected -= w.itemsPerPage
			}

		case tcell.KeyRight:
			if w.itemsPerPage > 0 && w.selected/w.itemsPerPage < w.pages-1 {
				newIndex := min(w.selected+w.itemsPerPage, len(w.filtered)-1)

				w.filtered[w.selected].widget = nil
				w.filtered[newIndex].widget = nil

				w.selected = newIndex
			}

		default:
		}
	}
}

func (w *List[T]) Draw(screen tcell.Screen, rect Rect) {
	selectedStrWidth := uniseg.StringWidth(w.SelectedStr)

	// Calculate paging
	w.itemsPerPage = rect.Height / w.ItemHeight
	if w.itemsPerPage == 0 {
		return
	}
	w.pages = divCeil(len(w.filtered), w.itemsPerPage)

	if w.pages > 1 {
		w.itemsPerPage = (rect.Height - 1) / w.ItemHeight
		if w.itemsPerPage == 0 {
			return
		}
		w.pages = divCeil(len(w.filtered), w.itemsPerPage)
	}

	page := w.selected / w.itemsPerPage
	offset := page * w.itemsPerPage

	// Draw items
	y := 0

	for i := offset; i < min(offset+w.itemsPerPage, len(w.filtered)); i++ {
		// Get item widget
		widget := w.filtered[i].widget

		if widget == nil {
			widget = w.ItemDisplayFn(w.items[w.filtered[i].index], i == w.selected && w.focused)
			w.filtered[i].widget = widget

			widget.CalcRequiredSize()
		}

		// Draw selected str
		if i == w.selected && w.focused {
			for i := range w.ItemHeight {
				screen.Put(rect.X, rect.Y+y+i, w.SelectedStr, w.SelectedStyle)
			}
		}

		// Draw item widget
		widget.Draw(screen, Rect{
			X:      rect.X + selectedStrWidth,
			Y:      rect.Y + y,
			Width:  max(rect.Width-selectedStrWidth, 0),
			Height: w.ItemHeight,
		})

		y += w.ItemHeight
	}

	// Draw page number
	if w.pages > 1 {
		style := Fg(color.Gray)

		str := fmt.Sprintf("page %d/%d", page+1, w.pages)
		strWidth := uniseg.StringWidth(str) + 1

		screen.PutStrStyled(rect.X+rect.Width-strWidth, rect.Y+rect.Height-1, str, style)
	}
}
