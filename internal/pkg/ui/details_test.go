package ui

import (
	td "github.com/goldenpathtechnologies/ci/testdata"
	"testing"
)

func Test_DetailsView_GetScrollAreaHandler(t *testing.T) {
	td.RunTextTestCases(func(data td.TextData, name string) {
		view := newDetailsView().SetText(data.Text)
		runScrollAreaTest(t, view, data.LongestLine, data.LineCount, name)
	})
}

func Test_DetailsView_GetScrollAreaHandler_WithWrap(t *testing.T) {
	td.RunTextTestCases(func(data td.TextData, name string) {
		view := newDetailsView().SetText(data.Text).SetWrap(true)
		view.SetRect(0, 0, data.LongestWrappedLine, data.LineCount)

		runScrollAreaTest(t, view, data.LongestWrappedLine, data.WrappedLineCount, name)
	})
}

func runScrollAreaTest(
	t *testing.T,
	view *DetailsView,
	expectedWidth int,
	expectedHeight int,
	description string,
) {
	width, height := view.handleScrollArea()

	if width != expectedWidth {
		t.Errorf(
			"%s: Expected scroll area width to be %d, got %d instead",
			description,
			expectedWidth,
			width)
	}

	if height != expectedHeight {
		t.Errorf(
			"%s: Expected scroll area height to be %d, got %d instead",
			description,
			expectedHeight,
			height)
	}
}

func Test_DetailsView_GetScrollPositionHandler(t *testing.T) {
	testText := td.TestText["LoremIpsum"].Text
	view := newDetailsView().SetText(testText)
	view.SetRect(0, 0, 10, 10)

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
		x, y := view.handleScrollPosition()
		if x != pos.x || y != pos.y {
			t.Errorf(
				"Expected scroll position to be (%d, %d), got (%d, %d) instead",
				pos.x, pos.y, x, y)
		}
	}
}

func Test_DetailsView_CreateDetailsPane(t *testing.T) {
	var d interface{} = CreateDetailsPane()

	_, isDetailsView := d.(*DetailsView)

	if !isDetailsView {
		t.Errorf("Expected type of object to be *DetailsView, got %T instead", d)
	}
}

func Test_DetailsView_GetText(t *testing.T) {
	td.RunTextTestCases(func(data td.TextData, name string) {
		view := newDetailsView()
		view.TextView.SetText(data.Text)
		result := view.GetText(true)
		runTextAccessTest(t, result, data.StrippedText, name)
	})
}

func Test_DetailsView_GetText_NoTagStripping(t *testing.T) {
	td.RunTextTestCases(func(data td.TextData, name string) {
		view := newDetailsView()
		view.TextView.SetText(data.Text)
		result := view.GetText(false)
		runTextAccessTest(t, result, data.Text, name)
	})
}

func Test_DetailsView_SetText(t *testing.T) {
	td.RunTextTestCases(func(data td.TextData, name string) {
		view := newDetailsView().SetText(data.Text)
		result := view.GetText(true)

		runTextAccessTest(t, result, data.StrippedText, name)
	})
}

func Test_DetailsView_SetText_ReturnsDetailsView(t *testing.T) {
	td.RunTextTestCases(func(data td.TextData, name string) {
		view := newDetailsView()
		viewAfterSet := view.SetText(data.Text)

		runReturnsDetailsViewTest(t, viewAfterSet, view, name)
	})
}

func runReturnsDetailsViewTest(
	t *testing.T,
	result *DetailsView,
	expected *DetailsView,
	description string,
) {
	if len(description) == 0 {
		description = "[Manual test]"
	}

	if expected != result {
		t.Errorf(
			"%s: Expected object '%v', got '%v' instead",
			description,
			expected,
			result)
	}
}

func runTextAccessTest(t *testing.T, result, expected string, description string) {
	if result != expected {
		t.Errorf(
			"%s: Expected text to be '%s', got '%s' instead",
			description,
			expected,
			result)
	}
}

func Test_DetailsView_SetText_LongestLine(t *testing.T) {
	td.RunTextTestCases(func(data td.TextData, name string) {
		view := newDetailsView()
		view.SetText(data.Text)

		runLineStatsTest(t, "LongestLine", view.LongestLine, data.LongestLine, name)
	})
}

func runLineStatsTest(t *testing.T, name string, result, expected int, description string) {
	if result != expected {
		t.Errorf(
			"%s: Expected %s to be %d, got %d instead",
			description,
			name,
			expected,
			result)
	}
}

func Test_DetailsView_SetText_LineCount(t *testing.T) {
	td.RunTextTestCases(func(data td.TextData, name string) {
		view := newDetailsView()
		view.SetText(data.Text)

		runLineStatsTest(t, "LineCount", view.LineCount, data.LineCount, name)
	})
}

func Test_DetailsView_SetWrap(t *testing.T) {
	view := newDetailsView()
	if view.HasWrap {
		t.Errorf("Expected HasWrap to be 'false', got 'true' instead")
	}

	view.SetWrap(true)
	if !view.HasWrap {
		t.Errorf("Expected HasWrap to be 'true', got 'false' instead")
	}
}

func Test_DetailsView_SetWrap_ReturnsDetailsView(t *testing.T) {
	view := newDetailsView()
	viewAfterWrap := view.SetWrap(true)

	runReturnsDetailsViewTest(t, viewAfterWrap, view, "")
}

func Test_DetailsView_SetWrap_LongestLine(t *testing.T) {
	td.RunTextTestCases(func(data td.TextData, name string) {
		view := newDetailsView()
		view.SetText(data.Text).SetRect(0, 0, data.LongestWrappedLine, 5)
		view.SetWrap(true)

		runLineStatsTest(t, "LongestLine", view.LongestLine, data.LongestWrappedLine, name)
	})
}

func Test_DetailsView_SetWrap_LineCount(t *testing.T) {
	td.RunTextTestCases(func(data td.TextData, name string) {
		view := newDetailsView()
		view.SetText(data.Text).SetRect(0, 0, data.LongestWrappedLine, 5)
		view.SetWrap(true)

		runLineStatsTest(t, "LineCount", view.LineCount, data.WrappedLineCount, name)
	})
}

func Test_DetailsView_SetWordWrap(t *testing.T) {
	view := newDetailsView()
	if view.HasWordWrap {
		t.Errorf("Expected HasWordWrap to be 'false', got 'true' instead")
	}

	view.SetWordWrap(true)
	if !view.HasWordWrap {
		t.Errorf("Expected HasWordWrap to be 'true', got 'false' instead")
	}
}

func Test_DetailsView_SetWordWrap_LongestLine(t *testing.T) {
	td.RunTextTestCases(func(data td.TextData, name string) {
		view := newDetailsView()
		view.SetText(data.Text).SetWrap(true).SetRect(0, 0, data.LongestWrappedLine, 5)
		view.SetWordWrap(true)

		runLineStatsTest(t, "LongestLine", view.LongestLine, data.LongestWordWrappedLine, name)
	})
}

func Test_DetailsView_SetWordWrap_LineCount(t *testing.T) {
	td.RunTextTestCases(func(data td.TextData, name string) {
		view := newDetailsView()
		view.SetText(data.Text).SetWrap(true).SetRect(0, 0, data.LongestWrappedLine, 5)
		view.SetWordWrap(true)

		runLineStatsTest(t, "LineCount", view.LineCount, data.WordWrappedLineCount, name)
	})
}

func Test_DetailsView_SetWordWrap_ReturnsDetailsView(t *testing.T) {
	view := newDetailsView()
	viewAfterWordWrap := view.SetWordWrap(true)

	runReturnsDetailsViewTest(t, viewAfterWordWrap, view, "")
}

func Test_DetailsView_SetRect(t *testing.T) {
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

func Test_DetailsView_SetRect_LongestLine(t *testing.T) {
	td.RunTextTestCases(func(data td.TextData, name string) {
		view := newDetailsView()
		view.
			SetText(data.Text).
			SetWrap(true).
			SetRect(0, 0, data.ViewWidth, 5)

		runLineStatsTest(t, "LongestLine", view.LongestLine, data.LongestWrappedLine, name)
	})
}

func Test_DetailsView_SetRect_LineCount(t *testing.T) {
	td.RunTextTestCases(func(data td.TextData, name string) {
		view := newDetailsView()
		view.
			SetText(data.Text).
			SetWrap(true).
			SetRect(0, 0, data.ViewWidth, 5)

		runLineStatsTest(t, "LineCount", view.LineCount, data.WrappedLineCount, name)
	})
}

func Test_DetailsView_SetRect_TextResize(t *testing.T) {
	td.RunTextTestCases(func(data td.TextData, name string) {
		view := newDetailsView().SetText(data.Text)
		textBeforeResize := view.GetText(false)

		view.SetRect(0, 0, 5, 5)
		textAfterResize := view.GetText(false)

		if textBeforeResize != textAfterResize {
			t.Errorf(
				"%s: Expected text '%s' after resize, got '%s' instead",
				name, textBeforeResize, textAfterResize)
		}
	})
}

func Test_DetailsView_calculateLineStats(t *testing.T) {
	view := newDetailsView()
	td.RunTextTestCases(func(data td.TextData, name string) {
		// Stripping tags if present since the function under test assumes verbatim text
		text := view.SetText(data.Text).GetText(true)
		testData := [][]lineStats{
			{
				calculateLineStats(text, data.ViewWidth, false, false),
				lineStats{data.LongestLine, data.LineCount},
			},
			{
				calculateLineStats(text, data.ViewWidth, true, false),
				lineStats{data.LongestWrappedLine, data.WrappedLineCount},
			},
			{
				calculateLineStats(text, data.ViewWidth, true, true),
				lineStats{data.LongestWordWrappedLine, data.WordWrappedLineCount},
			},
			{
				calculateLineStats(text, data.ViewWidth, false, true),
				lineStats{data.LongestLine, data.LineCount},
			},
		}

		runTest := func(t *testing.T, result, expected lineStats) {
			if result.longest != expected.longest || result.count != expected.count {
				t.Errorf(
					"%s: Expected longest line and line count to be (%d, %d), got (%d, %d) instead",
					name,
					expected.longest,
					expected.count,
					result.longest,
					result.count)
			}
		}

		for _, tData := range testData {
			runTest(t, tData[0], tData[1])
		}
	})
}
