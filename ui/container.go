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

func (w *Container) MaxSize() (int, int) {
	if w.maxWidth == 0 && w.maxHeight == 0 {
		maxWidth := 0
		maxHeight := 0

		if w.Direction == Vertical {
			for _, child := range w.Children {
				childMaxWidth, _ := child.MaxSize()
				maxWidth = max(maxWidth, childMaxWidth)
			}
		} else {
			for _, child := range w.Children {
				_, childMaxHeight := child.MaxSize()
				maxHeight = max(maxHeight, childMaxHeight)
			}
		}

		return maxWidth, maxHeight
	}

	return w.maxWidth, w.maxHeight
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
			maxWidth, maxHeight := child.MaxSize()

			if maxWidth > 0 && width > maxWidth {
				width = maxWidth
			}
			if maxHeight > 0 && height > maxHeight {
				height = maxHeight
			}

			w.requiredWidth = max(w.requiredWidth, width)
			w.requiredHeight += height
		}
	} else {
		for i, child := range w.Children {
			if i > 0 {
				w.requiredWidth += w.Gap
			}

			width, height := child.CalcRequiredSize()
			maxWidth, maxHeight := child.MaxSize()

			if maxWidth > 0 && width > maxWidth {
				width = maxWidth
			}
			if maxHeight > 0 && height > maxHeight {
				height = maxHeight
			}

			w.requiredWidth += width
			w.requiredHeight = max(w.requiredHeight, height)
		}
	}

	w.requiredWidth += (w.Padding + w.padding) * 2
	w.requiredHeight += (w.Padding + w.padding) * 2

	return w.requiredWidth, w.requiredHeight
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
			//    and which are fixed to their MaxSize.
			activeIndices := make([]int, 0, len(w.Children))
			for i := range w.Children {
				activeIndices = append(activeIndices, i)
			}

			// 3. Iteratively distribute space
			for len(activeIndices) > 0 {
				share := spaceToDistribute / len(activeIndices)
				remainder := spaceToDistribute % len(activeIndices)

				// We need to see if this share violates any MaxSize constraints.
				// We restart the distribution loop if we find ANY violation to ensure fairness.
				cappedFound := false

				// Identify indices that need to be capped
				// We iterate backwards to easily remove from activeIndices slice if needed,
				// or just build a new list for the next pass.
				nextActive := activeIndices[:0] // reusing storage

				for i, idx := range activeIndices {
					_, maxH := w.Children[idx].MaxSize()

					// Calculate what this child would get in this pass
					proposedHeight := share
					// (Optional) Distribute remainder to the first few items?
					// Usually simpler to apply remainder only in the final stable pass.
					// But strict checking might require looking at it now.
					// Let's assume remainder goes to the first ones in the list.
					if i < remainder {
						proposedHeight++
					}

					if maxH > 0 && proposedHeight > maxH {
						// Lock this child to maxH
						childHeights[idx] = maxH
						spaceToDistribute -= maxH
						cappedFound = true
					} else {
						// Keep active
						nextActive = append(nextActive, idx)
					}
				}

				if cappedFound {
					// A cap was hit, so the pool of space and active children changed.
					// Recalculate distribution for the remaining children.
					activeIndices = nextActive
				} else {
					// No caps hit, distribute remaining space to active children and finish.
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
			// Standard positioning (Start, Center, End)
			extraHeight := max(rect.Height-w.requiredHeight, 0)
			y += align(w.PrimaryAlignment, extraHeight)
		}

		for i, child := range w.Children {
			reqW, reqH := child.RequiredSize()
			maxW, maxH := child.MaxSize()

			if maxW > 0 && reqW > maxW {
				reqW = maxW
			}
			if maxH > 0 && reqH > maxH {
				reqH = maxH
			}

			// Determine Child Width (Cross Axis)
			childWidth := reqW
			childX := x

			if w.SecondaryAlignment == Stretch {
				childWidth = availWidth
				if maxW > 0 && childWidth > maxW {
					childWidth = maxW
				}
				childX = x
			} else {
				extraWidth := max(availWidth-reqW, 0)
				childX = x + align(w.SecondaryAlignment, extraWidth)
			}

			// Determine Child Height (Main Axis)
			childHeight := reqH

			if w.PrimaryAlignment == Stretch {
				// Use the calculated distributed height
				childHeight = childHeights[i]
				// Note: MaxSize is already baked into childHeights[i] by the solver above.
			}

			// Draw Child
			if y < rect.Y+rect.Height {
				rect := Rect{
					X:      childX,
					Y:      y,
					Width:  childWidth,
					Height: childHeight,
				}

				child.Draw(screen, rect)

				if onChildDraw != nil {
					onChildDraw(child, rect)
				}
			}

			y += childHeight + w.Gap
		}
	} else {
		// Horizontal Case - Mirror of Vertical Logic
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
					maxW, _ := w.Children[idx].MaxSize()

					proposedWidth := share
					if i < remainder {
						proposedWidth++
					}

					if maxW > 0 && proposedWidth > maxW {
						childWidths[idx] = maxW
						spaceToDistribute -= maxW
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
			maxW, maxH := child.MaxSize()

			if maxW > 0 && reqW > maxW {
				reqW = maxW
			}
			if maxH > 0 && reqH > maxH {
				reqH = maxH
			}

			// Determine Child Height (Cross Axis)
			childHeight := reqH
			childY := y

			if w.SecondaryAlignment == Stretch {
				childHeight = availHeight
				if maxH > 0 && childHeight > maxH {
					childHeight = maxH
				}
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

			// Draw Child
			if x < rect.X+rect.Width {
				rect := Rect{
					X:      x,
					Y:      childY,
					Width:  childWidth,
					Height: childHeight,
				}

				child.Draw(screen, rect)

				if onChildDraw != nil {
					onChildDraw(child, rect)
				}
			}

			x += childWidth + w.Gap
		}
	}
}
