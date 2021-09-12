package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"sync"
)

type Scrollable interface {
	tview.Primitive
	//SetBorder(show bool) *tview.Box
	//GetRect() (int, int, int, int)
	GetInnerRect() (int, int, int, int)
	//Draw(screen tcell.Screen)
	//HasFocus() bool
}

// TODO: Implement as many Primitive functions in this type as reasonable.
//  Draw(screen tcell.Screen)
//  GetRect() (int, int, int, int)
//  SetRect(x, y, width, height int)
//  InputHandler() func(event *tcell.EventKey, setFocus func(p Primitive))
//  Focus(delegate func(p Primitive))
//  HasFocus() bool
//  Blur() Note: Blur() isn't implemented by TextView but is by Box.
//  Note: Some of the above functions may not need implementing.
type ScrollView struct {
	sync.Mutex
	//*tview.TextView // TODO: Consider whether or not ScrollView should support more primitives.
	//*tview.Box
	Scrollable

	textColor tcell.Color
	backgroundColor tcell.Color
	//hasBorder bool
	contentSize func() (width, height int)
	scrollOffset func() (vScroll, hScroll int)
}

// TODO: If ScrollViews support more primitives, this function should have a *Box parameter.
func NewScrollView(s Scrollable) *ScrollView {
	return &ScrollView{
		//TextView: tview.NewTextView(),
		Scrollable: s,
		textColor: tview.Styles.PrimaryTextColor,
		//hasBorder: false,
		scrollOffset: func() (vScroll, hScroll int) {
			return 0, 0
		},
	}
}

////SetBorder sets the flag indicating whether the ScrollView should
////have a border.
//func (s *ScrollView) SetBorder(show bool) *ScrollView {
//	s.hasBorder = show
//	// TODO: Use s.Box if ScrollView will support more primitives
//	//s.TextView.SetBorder(show)
//	s.Scrollable.SetBorder(show)
//	return s
//}

//// TODO: Remove this function and let it be called by the underlying primitive
//// SetTextColor sets the color of the ScrollView text.
//func (s *ScrollView) SetTextColor(color tcell.Color) *ScrollView {
//	s.textColor = color
//	// TODO: Determine if the underlying Box needs to have its textColor changed as well
//	return s
//}

//// TODO: Remove this function and let it be called by the underlying primitive
//// SetBackgroundColor sets the color of the ScrollView background.
//func (s *ScrollView) SetBackgroundColor(color tcell.Color) *ScrollView {
//	s.backgroundColor = color
//	// TODO: Determine if the underlying Box needs to have its backgroundColor changed as well
//	return s
//}

// GetContentSize gets the size of the longest line and the number
// of lines of the ScrollView text, corresponding to the width
// and height of the scrollable area.
// TODO: The content size function should be a callback that is implemented
//  differently depending on what primitive is using the ScrollView.
func (s *ScrollView) GetContentSize() (width, height int) {
	return s.contentSize()
}

// TODO: The sizeFunction may be more appropriate as a parameter to NewScrollView.
func (s *ScrollView) SetContentSizeHandler(sizeFunction func() (width, height int)) *ScrollView {
	s.contentSize = sizeFunction
	return s
}

func (s *ScrollView) SetScrollOffsetHandler(handler func() (vScroll, hScroll int)) *ScrollView {
	s.scrollOffset = handler
	return s
}

func (s *ScrollView) GetScrollOffset() (vScroll, hScroll int) {
	return s.scrollOffset()
}

func (s *ScrollView) HasBorder() bool {
	x0, y0, height0, width0 := s.GetRect()
	x1, y1, height1, width1 := s.GetInnerRect()

	return x0 != x1 || y0 != y1 || height0-height1 >= 2 || width0-width1 >= 2
}

// TODO: Discovered the SetDrawFunc() function in tview.Box. I could somehow transform
//  this primitive to just a function that returns a function usable for SetDrawFunc().
//  One approach could be to use a thunk pattern, which may be a bit awkward, or I could
//  use handlers as parameters to the function that generates the Draw handler for
//  SetDrawFunc().
// Draw draws this primitive onto the screen.
func (s *ScrollView) Draw(screen tcell.Screen) {
	// TODO: Use s.Box if ScrollView will support more primitives
	//s.TextView.Draw(screen)
	s.Scrollable.Draw(screen)
	s.Lock()
	defer s.Unlock()

	x, y, width, height := s.GetRect()
	_, _, rectWidth, rectHeight := s.GetInnerRect()
	contentWidth, contentHeight := s.GetContentSize()

	//vScroll, hScroll := s.TextView.GetScrollOffset()
	vScroll, hScroll := s.GetScrollOffset()

	maxVScroll := contentHeight - rectHeight
	maxHScroll := contentWidth - rectWidth

	// For now, tview only supports border widths of 0 (no border) or 1
	var borderWidth int
	if s.HasBorder() {
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

func GetScrollBarDrawFunc(
	s Scrollable,
	contentHandler func() (width, height int),
	scrollHandler func() (vScroll, hScroll int),
	) func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {

	return func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {

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

		contentWidth, contentHeight := contentHandler()

		vScroll, hScroll := scrollHandler()

		maxVScroll := contentHeight - rectHeight
		maxHScroll := contentWidth - rectWidth

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

		return xI, yI, rectWidth, rectHeight
	}
}