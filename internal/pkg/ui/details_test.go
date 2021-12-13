package ui

import (
	"testing"
)

const (
	testText = `Lorem ipsum dolor sit amet, consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in
voluptate velit esse cillum dolore eu fugiat nulla pariatur.
Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia
deserunt mollit anim id est laborum.`
	longestLine = 74
	numLines = 7
	longestWrappedLine = 64
	viewPortWidth = longestLine - 10
	longestWordWrappedLine = 61
	numWrappedLines = 11
	numWordWrappedLines = 11
)

func TestGetScrollAreaHandler(t *testing.T) {
	view := newDetailsView().SetText(testText)
	runScrollAreaTest(t, view, longestLine, numLines)
}

func TestGetScrollAreaHandlerWithWrap(t *testing.T) {
	view := newDetailsView().SetText(testText).SetWrap(true)
	view.SetRect(0, 0, longestWrappedLine, numLines)

	runScrollAreaTest(t, view, viewPortWidth, numWrappedLines)
}

func runScrollAreaTest(
	t *testing.T,
	view *DetailsView,
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

func TestGetScrollPositionHandler(t *testing.T) {
	view := newDetailsView().SetText(testText)
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

func TestCreateDetailsPane(t *testing.T) {
	var d interface{} = CreateDetailsPane()

	_, isDetailsView := d.(*DetailsView)

	if !isDetailsView {
		t.Errorf("Expected type of object to be *DetailsView, got %T instead", d)
	}
}

func TestGetText(t *testing.T) {
	testData := []string{
		testText,
		`This is some test text
that has some newlines and
also one at the end
`,
		"This is some text that just has a space at the end ",
		// Note: Tabs are replaced with spaces in tview, therefore, this is not good test data
		//  https://github.com/rivo/tview/blob/2a6de950f73bdc70658f7e754d4b5593f15c8408/textview.go#L663
		// "This is some text with a tab at the end\t",
	}

	view := newDetailsView()

	for _, data := range testData {
		view.TextView.SetText(data)
		result := view.GetText(false)
		runTextAccessTest(t, result, data)
	}
}

func TestSetText(t *testing.T) {
	view := newDetailsView().SetText(testText)
	result := view.GetText(false)

	runTextAccessTest(t, result, testText)
}

func TestSetTextReturnsDetailsView(t *testing.T) {
	view := newDetailsView()
	viewAfterSet := view.SetText(testText)

	runReturnsDetailsViewTest(t, viewAfterSet, view)
}

func runReturnsDetailsViewTest(t *testing.T, result *DetailsView, expected *DetailsView) {
	if expected != result {
		t.Errorf("Expected object '%v', got '%v' instead", expected, result)
	}
}

func runTextAccessTest(t *testing.T, result, expected string) {
	if result != expected {
		t.Errorf(
			"Expected text to be '%s', got '%s' instead",
			expected,
			result)
	}
}

func TestSetTextLongestLine(t *testing.T) {
	view := newDetailsView()
	view.SetText(testText)

	runLineStatsTest(t, "LongestLine", view.LongestLine, longestLine)
}

func runLineStatsTest(t *testing.T, name string, result, expected int) {
	if result != expected {
		t.Errorf(
			"Expected %s to be %d, got %d instead",
			name,
			expected,
			result)
	}
}

func TestSetTextLineCount(t *testing.T) {
	view := newDetailsView()
	view.SetText(testText)

	runLineStatsTest(t, "LineCount", view.LineCount, numLines)
}

func TestSetWrap(t *testing.T) {
	view := newDetailsView()
	if view.HasWrap {
		t.Errorf("Expected HasWrap to be 'false', got 'true' instead")
	}

	view.SetWrap(true)
	if !view.HasWrap {
		t.Errorf("Expected HasWrap to be 'true', got 'false' instead")
	}
}

func TestSetWrapReturnsDetailsView(t *testing.T) {
	view := newDetailsView()
	viewAfterWrap := view.SetWrap(true)

	runReturnsDetailsViewTest(t, viewAfterWrap, view)
}

func TestSetWrapLongestLine(t *testing.T) {
	view := newDetailsView()
	view.SetText(testText).SetRect(0, 0, longestWrappedLine, 5)
	view.SetWrap(true)

	runLineStatsTest(t, "LongestLine", view.LongestLine, longestWrappedLine)
}

func TestSetWrapLineCount(t *testing.T) {
	view := newDetailsView()
	view.SetText(testText).SetRect(0, 0, longestWrappedLine, 5)
	view.SetWrap(true)

	runLineStatsTest(t, "LineCount", view.LineCount, numWrappedLines)
}

func TestSetWordWrap(t *testing.T) {
	view := newDetailsView()
	if view.HasWordWrap {
		t.Errorf("Expected HasWordWrap to be 'false', got 'true' instead")
	}

	view.SetWordWrap(true)
	if !view.HasWordWrap {
		t.Errorf("Expected HasWordWrap to be 'true', got 'false' instead")
	}
}

func TestSetWordWrapLongestLine(t *testing.T) {
	view := newDetailsView()
	view.SetText(testText).SetWrap(true).SetRect(0, 0, longestWrappedLine, 5)
	view.SetWordWrap(true)

	runLineStatsTest(t, "LongestLine", view.LongestLine, longestWordWrappedLine)
}

func TestSetWordWrapLineCount(t *testing.T) {
	view := newDetailsView()
	view.SetText(testText).SetWrap(true).SetRect(0, 0, longestWrappedLine, 5)
	view.SetWordWrap(true)

	runLineStatsTest(t, "LineCount", view.LineCount, numWordWrappedLines)
}

func TestSetWordWrapReturnsDetailsView(t *testing.T) {
	view := newDetailsView()
	viewAfterWordWrap := view.SetWordWrap(true)

	runReturnsDetailsViewTest(t, viewAfterWordWrap, view)
}

func TestSetRect(t *testing.T) {
	view := newDetailsView()
	view.SetRect(1, 2, 3, 4)

	x, y, width, height := view.TextView.GetRect()
	if !(x == 1 && y == 2 && width == 3 && height == 4) {
		t.Errorf(
			"Expected Rect to have dimensions (1, 2, 3, 4), got (%d, %d, %d, %d) instead",
			x, y, width, height)
	}

	x, y, width, height = view.TextView.GetInnerRect()
	if !(x == 1 && y == 2 && width == 3 && height == 4) {
		t.Errorf(
			"Expected InnerRect to have dimensions (1, 2, 3, 4), got (%d, %d, %d, %d) instead",
			x, y, width, height)
	}
}

func TestSetRectLongestLine(t *testing.T) {
	view := newDetailsView()
	view.
		SetText("aaaaaaaa").
		SetWrap(true).
		SetRect(0, 0, 2, 5)

	runLineStatsTest(t, "LongestLine", view.LongestLine, 2)
}

func TestSetRectLineCount(t *testing.T) {
	view := newDetailsView()
	view.
		SetText("bbbbbbbbbbbbbbbb").
		SetWrap(true).
		SetRect(0, 0, 8, 5)

	runLineStatsTest(t, "LineCount", view.LineCount, 2)
}

func TestSetRectTextResize(t *testing.T) {
	view := newDetailsView().SetText(testText)
	textBeforeResize := view.GetText(false)

	view.SetRect(0, 0, 5, 5)
	textAfterResize := view.GetText(false)

	if textBeforeResize != textAfterResize {
		t.Errorf(
			"Expected text '%s' after resize, got '%s' instead",
			textBeforeResize, textAfterResize)
	}
}

func TestCalculateLineStats(t *testing.T) {
	testData := [][]lineStats{
		{
			calculateLineStats(testText, longestLine, false, false),
			lineStats{longestLine, numLines},
		},
		{
			calculateLineStats(testText, longestWrappedLine, true, false),
			lineStats{longestWrappedLine, numWrappedLines},
		},
		{
			calculateLineStats(testText, longestWrappedLine, true, true),
			lineStats{longestWordWrappedLine, numWordWrappedLines},
		},
		{
			calculateLineStats(testText, longestLine, false, true),
			lineStats{longestLine, numLines},
		},
		{
			calculateLineStats("bbbbbbbbbbbbbbbb", 8, true, false),
			lineStats{8, 2},
		},
		{
			calculateLineStats("hhhhhhhhhhhhhhhhhhhhhhhh", 6, true, false),
			lineStats{6, 4},
		},
		{
			calculateLineStats("aaaa\naaaa\naaaa\naaaa", 2, true, false),
			lineStats{2, 8},
		},
		{
			calculateLineStats("aaaaa\naaaa\naaaa", 2, true, false),
			lineStats{2, 7},
		},
	}

	runTest := func(t *testing.T, result, expected lineStats) {
		if result.longest != expected.longest || result.count != expected.count {
			t.Errorf(
				"Expected longest line and line count to be (%d, %d), got (%d, %d) instead",
				expected.longest,
				expected.count,
				result.longest,
				result.count)
		}
	}

	for _, data := range testData {
		runTest(t, data[0], data[1])
	}
}
