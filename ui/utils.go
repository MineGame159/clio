package ui

import (
	"strings"

	"github.com/gdamore/tcell/v3"
)

// Direction

type Direction uint8

const (
	Vertical Direction = iota
	Horizontal
)

// Alignment

type Alignment uint8

const (
	Start Alignment = iota
	Center
	End
	Stretch
)

// Rect

type Rect struct {
	X, Y          int
	Width, Height int
}

func (r Rect) Pad(padding int) Rect {
	if r.Width < padding*2 || r.Height < padding*2 {
		return Rect{
			X:      r.X,
			Y:      r.Y,
			Width:  0,
			Height: 0,
		}
	}

	return Rect{
		X:      r.X + padding,
		Y:      r.Y + padding,
		Width:  r.Width - padding*2,
		Height: r.Height - padding*2,
	}
}

// Helpers

func Fg(color tcell.Color) tcell.Style {
	return tcell.StyleDefault.Foreground(color)
}

func LeftMargin(widget Widget, value int) Widget {
	p := NewParagraph(strings.Repeat(" ", value))
	p.SetMaxWidth(value)

	return &Container{
		Direction:          Horizontal,
		PrimaryAlignment:   Stretch,
		SecondaryAlignment: Stretch,
		Children:           []Widget{p, widget},
	}
}

func FilterFn[T any](query string, filteredValueFn func(item T) string) func(item T) bool {
	words := strings.Split(strings.ToLower(query), " ")

	return func(item T) bool {
		filteredValue := strings.ToLower(filteredValueFn(item))

		for _, word := range words {
			if !strings.Contains(filteredValue, word) {
				return false
			}
		}

		return true
	}
}

func align(alignment Alignment, additionalParentSize int) int {
	switch alignment {
	case Start, Stretch:
		return 0
	case Center:
		return additionalParentSize / 2
	case End:
		return additionalParentSize
	default:
		panic("Invalid Alignment")
	}
}

func divCeil(x, y int) int {
	return 1 + ((x - 1) / y)
}

func putStr(screen tcell.Screen, x, y, maxX int, str string, style tcell.Style) int {
	localX := 0

	for str != "" && x+localX < maxX {
		var w int
		str, w = screen.Put(x+localX, y, str, style)

		if w == 0 {
			break
		}

		localX += w
	}

	return localX
}
