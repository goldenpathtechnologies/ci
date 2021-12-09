package ui

import (
	"github.com/rivo/tview"
	"testing"
)

const (
	text = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, 
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. 
Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut 
aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in 
voluptate velit esse cillum dolore eu fugiat nulla pariatur. 
Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia 
deserunt mollit anim id est laborum.`
	longestLine = 75
	numLines = 7
	longestWrappedLine = 65
	numWrappedLines = 11
)

func runScrollAreaTest(
	t *testing.T,
	view *tview.TextView,
	expectedWidth int,
	expectedHeight int,
) {
	handler := getScrollAreaHandler(view)

	width, height := handler()

	if width != expectedWidth {
		t.Errorf("Expected scroll area width to be %d, got %d instead", expectedWidth, width)
	}

	if height != expectedHeight {
		t.Errorf("Expected scroll area height to be %d, got %d instead", expectedHeight, height)
	}
}

func TestGetScrollAreaHandler(t *testing.T) {
	view := tview.NewTextView().SetText(text)

	runScrollAreaTest(t, view, longestLine, numLines)
}

func TestGetScrollAreaHandlerWithWrap(t *testing.T) {
	view := tview.NewTextView().SetText(text).SetWrap(true)
	view.SetRect(0, 0, longestWrappedLine, numLines)

	runScrollAreaTest(t, view, longestLine-10, numWrappedLines)
}

func TestGetScrollPositionHandler(t *testing.T) {
	view := tview.NewTextView().SetText(text)
	view.SetRect(0, 0, 10, 10)

	handler := getScrollPositionHandler(view)

	scrollData := []struct {
		x int
		y int
	}{
		{1, 1},
		{3, 5},
		{2, 7},
		{10, 10},
	}

	for _, pos := range scrollData {
		view.ScrollTo(pos.x, pos.y)
		x, y := handler()
		if x != pos.x || y != pos.y {
			t.Errorf(
				"Expected scroll position to be (%d, %d), got (%d, %d) instead",
				pos.x, pos.y, x, y)
		}
	}
}
