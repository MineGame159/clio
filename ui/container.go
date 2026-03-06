package ui

import "github.com/gdamore/tcell/v3"

type Container struct {
	baseWidget

	Direction          Direction
	PrimaryAlignment   Alignment
	SecondaryAlignment Alignment

	Padding int
	Gap     int

	padding int

	Children []Widget
}

func (w *Container) CalcRequiredSize() (int, int) {
	w.requiredWidth = 0
	w.requiredHeight = 0

	if w.Direction == Vertical {
		for i, child := range w.Children {
			if i > 0 {
				w.requiredHeight += w.Gap
			}

			width, height := child.CalcRequiredSize()
			width, height = child.LimitSize(width, height)

			w.requiredWidth = max(w.requiredWidth, width)
			w.requiredHeight += height
		}
	} else {
		for i, child := range w.Children {
			if i > 0 {
				w.requiredWidth += w.Gap
			}

			width, height := child.CalcRequiredSize()
			width, height = child.LimitSize(width, height)

			w.requiredWidth += width
			w.requiredHeight = max(w.requiredHeight, height)
		}
	}

	w.requiredWidth += (w.Padding + w.padding) * 2
	w.requiredHeight += (w.Padding + w.padding) * 2

	return w.requiredWidth, w.requiredHeight
}

func (w *Container) LimitSize(width, height int) (int, int) {
	if w.Direction == Vertical {
		for _, child := range w.Children {
			width, _ = child.LimitSize(width, 0)
		}
	} else {
		for _, child := range w.Children {
			_, height = child.LimitSize(0, height)
		}
	}

	return width, height
}

func (w *Container) HandleEvent(event any) {
	for _, child := range w.Children {
		child.HandleEvent(event)
	}
}

func (w *Container) Draw(screen tcell.Screen, rect Rect) {
	w.draw(screen, rect, nil)
}

func (w *Container) draw(screen tcell.Screen, rect Rect, onChildDraw func(child Widget, rect Rect)) {
	padding := w.Padding + w.padding

	availWidth := max(rect.Width-padding*2, 0)
	availHeight := max(rect.Height-padding*2, 0)

	x := rect.X + padding
	y := rect.Y + padding

	if w.Direction == Vertical {
		// Prepare a slice to store calculated heights for every child
		childHeights := make([]int, len(w.Children))

		if w.PrimaryAlignment == Stretch {
			// 1. Calculate total available space for children (minus gaps)
			totalGap := w.Gap * (len(w.Children) - 1)
			spaceToDistribute := max(availHeight-totalGap, 0)

			// 2. Track which children are still taking part in the distribution
			activeIndices := make([]int, 0, len(w.Children))
			for i := range w.Children {
				activeIndices = append(activeIndices, i)
			}

			// 3. Iteratively distribute space
			for len(activeIndices) > 0 {
				share := spaceToDistribute / len(activeIndices)
				remainder := spaceToDistribute % len(activeIndices)

				cappedFound := false
				nextActive := activeIndices[:0] // reusing storage

				for i, idx := range activeIndices {
					proposedHeight := share
					if i < remainder {
						proposedHeight++
					}

					// Use available Width if stretched to allow aspect ratio constraints to work
					crossWidth := 0
					if w.SecondaryAlignment == Stretch {
						crossWidth = availWidth
					}

					_, limitedHeight := w.Children[idx].LimitSize(crossWidth, proposedHeight)

					if limitedHeight < proposedHeight {
						childHeights[idx] = limitedHeight
						spaceToDistribute -= limitedHeight
						cappedFound = true
					} else {
						nextActive = append(nextActive, idx)
					}
				}

				if cappedFound {
					activeIndices = nextActive
				} else {
					for i, idx := range activeIndices {
						h := share
						if i < remainder {
							h++
						}
						childHeights[idx] = h
					}
					break
				}
			}
		} else {
			extraHeight := max(rect.Height-w.requiredHeight, 0)
			y += align(w.PrimaryAlignment, extraHeight)
		}

		for i, child := range w.Children {
			reqW, reqH := child.RequiredSize()
			reqW, reqH = child.LimitSize(reqW, reqH)

			// Determine Child Width (Cross Axis)
			childWidth := reqW
			childX := x

			if w.SecondaryAlignment == Stretch {
				childWidth = availWidth
				childX = x
			} else {
				extraWidth := max(availWidth-reqW, 0)
				childX = x + align(w.SecondaryAlignment, extraWidth)
			}

			// Determine Child Height (Main Axis)
			childHeight := reqH

			if w.PrimaryAlignment == Stretch {
				childHeight = childHeights[i]
			}

			// Apply constraints with both dimensions known
			childWidth, childHeight = child.LimitSize(childWidth, childHeight)

			if y < rect.Y+rect.Height {
				r := Rect{
					X:      childX,
					Y:      y,
					Width:  childWidth,
					Height: childHeight,
				}

				child.Draw(screen, r)

				if onChildDraw != nil {
					onChildDraw(child, r)
				}
			}

			y += childHeight + w.Gap
		}
	} else {
		// Horizontal Case
		childWidths := make([]int, len(w.Children))

		if w.PrimaryAlignment == Stretch {
			totalGap := w.Gap * (len(w.Children) - 1)
			spaceToDistribute := max(availWidth-totalGap, 0)

			activeIndices := make([]int, 0, len(w.Children))
			for i := range w.Children {
				activeIndices = append(activeIndices, i)
			}

			for len(activeIndices) > 0 {
				share := spaceToDistribute / len(activeIndices)
				remainder := spaceToDistribute % len(activeIndices)

				cappedFound := false
				nextActive := activeIndices[:0]

				for i, idx := range activeIndices {
					proposedWidth := share
					if i < remainder {
						proposedWidth++
					}

					// Use available Height if stretched to allow aspect ratio constraints to work
					crossHeight := 0
					if w.SecondaryAlignment == Stretch {
						crossHeight = availHeight
					}

					limitedWidth, _ := w.Children[idx].LimitSize(proposedWidth, crossHeight)

					if limitedWidth < proposedWidth {
						childWidths[idx] = limitedWidth
						spaceToDistribute -= limitedWidth
						cappedFound = true
					} else {
						nextActive = append(nextActive, idx)
					}
				}

				if cappedFound {
					activeIndices = nextActive
				} else {
					for i, idx := range activeIndices {
						w := share
						if i < remainder {
							w++
						}
						childWidths[idx] = w
					}
					break
				}
			}
		} else {
			extraWidth := max(rect.Width-w.requiredWidth, 0)
			x += align(w.PrimaryAlignment, extraWidth)
		}

		for i, child := range w.Children {
			reqW, reqH := child.RequiredSize()
			reqW, reqH = child.LimitSize(reqW, reqH)

			// Determine Child Height (Cross Axis)
			childHeight := reqH
			childY := y

			if w.SecondaryAlignment == Stretch {
				childHeight = availHeight
				childY = y
			} else {
				extraHeight := max(availHeight-reqH, 0)
				childY = y + align(w.SecondaryAlignment, extraHeight)
			}

			// Determine Child Width (Main Axis)
			childWidth := reqW

			if w.PrimaryAlignment == Stretch {
				childWidth = childWidths[i]
			}

			// Apply constraints with both dimensions known
			childWidth, childHeight = child.LimitSize(childWidth, childHeight)

			if x < rect.X+rect.Width {
				r := Rect{
					X:      x,
					Y:      childY,
					Width:  childWidth,
					Height: childHeight,
				}

				child.Draw(screen, r)

				if onChildDraw != nil {
					onChildDraw(child, r)
				}
			}

			x += childWidth + w.Gap
		}
	}
}
