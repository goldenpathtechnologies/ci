package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
	"sync"
)

type ScrollView struct {
	sync.Mutex
	*tview.TextView

	textColor tcell.Color
	backgroundColor tcell.Color
}

func NewScrollView() *ScrollView {
	return &ScrollView{
		TextView: tview.NewTextView(),
		textColor: tview.Styles.PrimaryTextColor,
	}
}

func (s *ScrollView) SetTextColor(color tcell.Color) *ScrollView {
	s.textColor = color
	return s
}

func (s *ScrollView) SetBackgroundColor(color tcell.Color) *ScrollView {
	s.backgroundColor = color
	return s
}

func (s *ScrollView) GetContentSize() (width, height int) {
	text := s.GetText(true)
	// TODO: This breaks when word wrap is enabled in the TextView as wrapped
	//  lines are not delimited by a new line externally. Ideally, I'd just
	//  access the s.longestLine and s.pageSize fields, but those are
	//  not exported from tview and therefore inaccessible.
	lines := strings.Split(text, "\n")
	longestLine := ""

	for _, v := range lines {
		if len(v) > len(longestLine) {
			longestLine = v
		}
	}

	return len(longestLine), len(lines)
}

func (s *ScrollView) Draw(screen tcell.Screen) {
	s.TextView.DrawForSubclass(screen, s)
	s.Lock()
	defer s.Unlock()

	s.TextView.Draw(screen)

	vOffset, hOffset := s.TextView.GetScrollOffset()

	x, y, width, height := s.GetRect()
	_, _, rectWidth, rectHeight := s.GetInnerRect()

	contentWidth, contentHeight := s.GetContentSize()

	scrollHeight := height - 2 // TODO: Factor in whether box border exists instead of hardcode of 2.
	scrollWidth := width - 2 // TODO: Same as above.

	vThumbSize := int(float32(rectHeight) / float32(contentHeight) * float32(scrollHeight))
	hThumbSize := int(float32(rectWidth) / float32(contentWidth) * float32(scrollWidth))

	maxVThumbOffset := scrollHeight - vThumbSize
	maxHThumbOffset := scrollWidth - hThumbSize

	maxVOffset := contentHeight - rectHeight
	maxHOffset := contentWidth - rectWidth

	vThumbOffset := int(float32(maxVThumbOffset) * float32(vOffset) / float32(maxVOffset))
	hThumbOffset := int(float32(maxHThumbOffset) * float32(hOffset) / float32(maxHOffset))

	if vThumbOffset == 0 && vOffset > 0 {
		vThumbOffset = 1
	}

	if hThumbOffset == 0 && hOffset > 0 {
		hThumbOffset = 1
	}

	scrollBarStyle := tcell.StyleDefault.Foreground(tcell.ColorLightGray)
	thumbStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow)

	for i := y + 1; i < height - 1; i++ {
		// TODO: the +1 below is taking the border into consideration, fyi.
		if i >= y + vThumbOffset + 1 && i <= y + vThumbOffset + vThumbSize {
			screen.SetContent(x + width - 1, i, tcell.RuneVLine, nil, thumbStyle)
		} else {
			screen.SetContent(x + width - 1, i, tcell.RuneVLine, nil, scrollBarStyle)
		}
	}

	for j := x + 1; j < width - 1; j++ {
		// TODO: the +1 below is taking the border into consideration, fyi.
		if j >= x + hThumbOffset + 1 && j <= x + hThumbOffset + hThumbSize {
			screen.SetContent(j, y + height - 1, tcell.RuneHLine, nil, thumbStyle)
		} else {
			screen.SetContent(j, y + height - 1, tcell.RuneHLine, nil, scrollBarStyle)
		}
	}
}