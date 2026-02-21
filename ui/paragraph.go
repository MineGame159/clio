package ui

import (
	"strings"

	"github.com/gdamore/tcell/v3"
	"github.com/rivo/uniseg"
)

type Paragraph struct {
	baseWidget

	Spans     []Span
	Alignment Alignment

	lines [][]Span
}

type Span struct {
	Text  string
	Style tcell.Style
}

func NewParagraph(text string) *Paragraph {
	return &Paragraph{Spans: []Span{{text, tcell.StyleDefault}}}
}

func NewParagraphStyled(text string, style tcell.Style) *Paragraph {
	return &Paragraph{Spans: []Span{{text, style}}}
}

func (w *Paragraph) Add(text string) {
	w.Spans = append(w.Spans, Span{text, tcell.StyleDefault})
}

func (w *Paragraph) AddStyled(text string, style tcell.Style) {
	w.Spans = append(w.Spans, Span{text, style})
}

func (w *Paragraph) CalcRequiredSize() (int, int) {
	// Split spans into lines
	w.lines = nil

	var line []Span

	for _, span := range w.Spans {
		splits := strings.Split(span.Text, "\n")

		for i, split := range splits {
			if i > 0 && len(line) > 0 {
				w.lines = append(w.lines, line)
				line = nil
			}

			line = append(line, Span{split, span.Style})
		}
	}

	if len(line) > 0 {
		w.lines = append(w.lines, line)
	}

	// Calculate required size
	w.requiredWidth = 0
	w.requiredHeight = len(w.lines)

	for _, line := range w.lines {
		width := 0

		for _, span := range line {
			width += uniseg.StringWidth(span.Text)
		}

		w.requiredWidth = max(w.requiredWidth, width)
	}

	return w.requiredWidth, w.requiredHeight
}

func (w *Paragraph) HandleEvent(_ any) {
}

func (w *Paragraph) Draw(screen tcell.Screen, rect Rect) {
	for y, line := range w.lines {
		if y >= rect.Height {
			break
		}

		// Calculate line width
		lineWidth := 0

		for _, span := range line {
			lineWidth += uniseg.StringWidth(span.Text)
		}

		// Draw spans
		x := align(w.Alignment, max(w.requiredWidth-lineWidth, 0))

		for _, span := range line {
			x += putStr(screen, rect.X+x, rect.Y+y, rect.Width, span.Text, span.Style)
		}
	}
}
