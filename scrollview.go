package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
	"sync"
)

type ScrollView struct {
	sync.Mutex
	*tview.TextView // TODO: Consider whether or not ScrollView should support more primitives.

	textColor tcell.Color
	backgroundColor tcell.Color
	hasBorder bool
}

// TODO: If ScrollViews support more primitives, this function should have a *Box parameter.
func NewScrollView() *ScrollView {
	return &ScrollView{
		TextView: tview.NewTextView(),
		textColor: tview.Styles.PrimaryTextColor,
		hasBorder: false,
	}
}

// SetBorder sets the flag indicating whether the ScrollView should
// have a border.
func (s *ScrollView) SetBorder(show bool) *ScrollView {
	s.hasBorder = show
	// TODO: Use s.Box if ScrollView will support more primitives
	s.TextView.SetBorder(show)
	return s
}

// SetTextColor sets the color of the ScrollView text.
func (s *ScrollView) SetTextColor(color tcell.Color) *ScrollView {
	s.textColor = color
	return s
}

// SetBackgroundColor sets the color of the ScrollView background.
func (s *ScrollView) SetBackgroundColor(color tcell.Color) *ScrollView {
	s.backgroundColor = color
	return s
}

// GetContentSize gets the size of the longest line and the number
// of lines of the ScrollView text, corresponding to the width
// and height of the scrollable area.
// TODO: The content size function should be a callback that is implemented
//  differently depending on what primitive is using the ScrollView.
func (s *ScrollView) GetContentSize() (width, height int) {
	text := s.GetText(true)
	// TODO: This breaks when word wrap is enabled in the TextView as wrapped
	//  lines are not delimited by a new line externally. Ideally, I'd just
	//  access the s.longestLine and s.pageSize fields, but those are
	//  not exported from tview.TextView and therefore inaccessible.
	lines := strings.Split(text, "\n")
	longestLine := ""

	for _, v := range lines {
		if len(v) > len(longestLine) {
			longestLine = v
		}
	}

	return len(longestLine), len(lines)
}

// Draw draws this primitive onto the screen.
func (s *ScrollView) Draw(screen tcell.Screen) {
	// TODO: Use s.Box if ScrollView will support more primitives
	s.TextView.Draw(screen)
	s.Lock()
	defer s.Unlock()

	x, y, width, height := s.GetRect()
	_, _, rectWidth, rectHeight := s.GetInnerRect()
	contentWidth, contentHeight := s.GetContentSize()

	vScroll, hScroll := s.TextView.GetScrollOffset()

	maxVScroll := contentHeight - rectHeight
	maxHScroll := contentWidth - rectWidth

	// For now, tview only supports border widths of 0 (no border) or 1
	var borderWidth int
	if s.hasBorder {
		borderWidth = 1
	} else {
		borderWidth = 0
	}

	vScrollBarSize := height - (borderWidth * 2)
	hScrollBarSize := width - (borderWidth * 2)

	vThumbSize := int(float32(rectHeight) / float32(contentHeight) * float32(vScrollBarSize))
	hThumbSize := int(float32(rectWidth) / float32(contentWidth) * float32(hScrollBarSize))

	maxVThumbScroll := vScrollBarSize - vThumbSize
	maxHThumbScroll := hScrollBarSize - hThumbSize

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

	// TODO: Do not draw horizontal scrollbar if the underlying TextView enables wrap
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
}
