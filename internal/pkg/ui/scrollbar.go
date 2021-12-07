package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Scrollable interface {
	tview.Primitive
	GetInnerRect() (int, int, int, int)
}

func GetScrollBarDrawFunc(
	s Scrollable,
	getScrollArea func() (width, height int),
	getScrollPosition func() (vScroll, hScroll int),
	) func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {

	return func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {

		// Note: Currently, the code is unable to distinguish between a Scrollable that has
		//  a border and one that doesn't have a border but instead has padding.
		xI, yI, rectWidth, rectHeight := s.GetInnerRect()
		hasBorder := func() bool {
			return x != xI || y != yI || height-rectHeight >= 2 || width-rectWidth >= 2
		}

		// For now, tview only supports border widths of 0 (no border) or 1
		var borderWidth int
		if hasBorder() {
			borderWidth = 1
		} else {
			borderWidth = 0
		}

		contentWidth, contentHeight := getScrollArea()

		vScroll, hScroll := getScrollPosition()

		// Maximum scrollable distance
		maxVScroll := contentHeight - rectHeight
		maxHScroll := contentWidth - rectWidth

		vScrollBarSize := height - (borderWidth * 2)
		hScrollBarSize := width - (borderWidth * 2)

		vThumbSize := int(float32(rectHeight) / float32(contentHeight) * float32(vScrollBarSize))
		hThumbSize := int(float32(rectWidth) / float32(contentWidth) * float32(hScrollBarSize))

		// Maximum scrollable distance for the scrollbar thumb
		maxVThumbScroll := vScrollBarSize - vThumbSize
		maxHThumbScroll := hScrollBarSize - hThumbSize

		// Position of scrollbar thumb
		vThumbScroll := int(float32(maxVThumbScroll) * float32(vScroll) / float32(maxVScroll))
		hThumbScroll := int(float32(maxHThumbScroll) * float32(hScroll) / float32(maxHScroll))

		// Ensure that the scrollbar thumb is always offset when the content has scrolled.
		if vThumbScroll == 0 && vScroll > 0 {
			vThumbScroll = 1
		}

		// Same as above but in the horizontal direction.
		if hThumbScroll == 0 && hScroll > 0 {
			hThumbScroll = 1
		}

		scrollBarStyle := tcell.StyleDefault.Foreground(tcell.ColorLightGray)
		thumbStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow)

		var vThumbRune, hThumbRune rune
		if s.HasFocus() {
			vThumbRune, hThumbRune = tview.Borders.VerticalFocus, tview.Borders.HorizontalFocus
		} else {
			vThumbRune, hThumbRune = tview.Borders.Vertical, tview.Borders.Horizontal
		}

		// Draw the vertical scrollbar
		if contentHeight > rectHeight {
			scrollY := y+borderWidth
			for i := scrollY; i < scrollY+vScrollBarSize; i++ {
				if i >= scrollY+vThumbScroll && i < scrollY+vThumbScroll+vThumbSize {
					screen.SetContent(x+width-1, i, vThumbRune, nil, thumbStyle)
				} else {
					screen.SetContent(x+width-1, i, tview.Borders.Vertical, nil, scrollBarStyle)
				}
			}
		}

		// Draw the horizontal scrollbar
		if contentWidth > rectWidth {
			scrollX := x+borderWidth
			for j := scrollX; j < scrollX+hScrollBarSize; j++ {
				if j >= scrollX+hThumbScroll && j < scrollX+hThumbScroll+hThumbSize {
					screen.SetContent(j, y+height-1, hThumbRune, nil, thumbStyle)
				} else {
					screen.SetContent(j, y+height-1, tview.Borders.Horizontal, nil, scrollBarStyle)
				}
			}
		}

		return xI, yI, rectWidth, rectHeight
	}
}