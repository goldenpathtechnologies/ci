package ui

import (
	"github.com/rivo/tview"
	"strings"
)

func CreateDetailsPane() *tview.TextView {
	details := tview.NewTextView()

	details.
		SetDynamicColors(true).
		SetWrap(false).
		SetTitle("Details").
		SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetDrawFunc(GetScrollBarDrawFunc(
			details,
			getScrollAreaHandler(details),
			getScrollPositionHandler(details)))

	return details
}

func getScrollAreaHandler(view *tview.TextView) func() (width, height int) {
	return func() (width, height int) {
		text := view.GetText(true)
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
}

func getScrollPositionHandler(view *tview.TextView) func() (vScroll, hScroll int) {
	return func() (vScroll, hScroll int) {
		return view.GetScrollOffset()
	}
}
