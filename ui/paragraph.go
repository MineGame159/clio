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
	y := 0

	for _, logicalLine := range w.lines {
		// Measure line width
		totalLineWidth := calcLineWidth(logicalLine)

		// Skip word wrapping if the line fits inside rectangle
		if totalLineWidth <= rect.Width {
			if w.drawLine(screen, rect, &y, totalLineWidth, logicalLine) {
				return
			}

			continue
		}

		// Word wrap
		var currentVisualLine []Span
		currentLineWidth := 0

		for _, span := range logicalLine {
			i := 0

			for word := range strings.SplitSeq(span.Text, " ") {
				wordWidth := uniseg.StringWidth(word)

				spaceWidth := 0
				prefix := ""

				if i > 0 {
					spaceWidth = 1
					prefix = " "
				}

				if currentLineWidth+spaceWidth+wordWidth > rect.Width && len(currentVisualLine) > 0 {
					if w.drawLine(screen, rect, &y, currentLineWidth, currentVisualLine) {
						return
					}

					currentVisualLine = currentVisualLine[:0]
					currentLineWidth = 0
					prefix = ""
					spaceWidth = 0
				}

				currentVisualLine = append(currentVisualLine, Span{
					Text:  prefix + word,
					Style: span.Style,
				})

				currentLineWidth += spaceWidth + wordWidth
				i++
			}
		}

		if len(currentVisualLine) > 0 {
			if w.drawLine(screen, rect, &y, currentLineWidth, currentVisualLine) {
				return
			}
		}
	}
}

func (w *Paragraph) drawLine(screen tcell.Screen, rect Rect, y *int, lineWidth int, spans []Span) bool {
	if *y >= rect.Height {
		return true
	}

	x := rect.X + align(w.Alignment, max(rect.Width-lineWidth, 0))

	for _, span := range spans {
		x += putStr(screen, x, rect.Y+*y, rect.X+rect.Width, span.Text, span.Style)
	}

	*y++
	return false
}

func calcLineWidth(spans []Span) int {
	width := 0

	for _, span := range spans {
		width += uniseg.StringWidth(span.Text)
	}

	return width
}
