package ui

import (
	"github.com/rivo/tview"
	"strings"
	"unicode/utf8"
)

const detailsViewTitle = "Details"

// DetailsView is a wrapper for tview.TextView with better support for scrolling
// content. Overridden functions of this struct must be called before any inherited
// functions from tview.TextView.
type DetailsView struct {
	*tview.TextView
	LongestLine int
	LineCount   int
	HasWrap     bool
	HasWordWrap bool
}

// lineStats is an internal type that keeps track of the longest line of the DetailsView
// content (width) and the line count (height).
type lineStats struct {
	longest int
	count   int
}

// CreateDetailsView creates a new instance of DetailsView and initializes it with default
// settings.
func CreateDetailsView() *DetailsView {
	details := &DetailsView{
		TextView:    tview.NewTextView().SetDynamicColors(true),
		LongestLine: 0,
		LineCount:   0,
		HasWrap:     false,
		HasWordWrap: false,
	}

	details.
		SetWrap(false).
		SetTitle(detailsViewTitle).
		SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetDrawFunc(GetScrollBarDrawFunc(
			details,
			details.handleScrollArea,
			details.handleScrollPosition))

	return details
}

// SetWrap sets the wrap setting in the underlying tview.TextView and recalculates the content
// width and height.
func (d *DetailsView) SetWrap(wrap bool) *DetailsView {
	d.HasWrap = wrap
	d.TextView.SetWrap(wrap)
	return d.refreshLineStats()
}

// refreshLineStats calculates the width and height of the scrollable content of the DetailsView.
func (d *DetailsView) refreshLineStats() *DetailsView {
	_, _, viewWidth, _ := d.TextView.GetInnerRect()
	lineData := calculateLineStats(d.GetText(true), viewWidth, d.HasWrap, d.HasWordWrap)
	d.LongestLine = lineData.longest
	d.LineCount = lineData.count
	return d
}

// calculateLineStats calculates the width and height of a scrollable area given the supplied text,
// the maximum width, and whether the content wraps normally, on words, or not at all.
func calculateLineStats(text string, maxWidth int, wrap, wordWrap bool) lineStats {
	if len(text) == 0 {
		return lineStats{0, 1}
	}

	lines := strings.Split(text, "\n")

	if wrap && wordWrap {
		lines = tview.WordWrap(text, maxWidth)

		// Note: tview.WordWrap() trims trailing newlines.
		if text[len(text)-1] == '\n' {
			lines = append(lines, "")
		}
	}

	longestLine := 0
	wrappedLines := 0

	for _, line := range lines {
		if len(line) > longestLine {
			longestLine = len(line)
		}

		if wrap && !wordWrap && len(line) > maxWidth {
			wrappedLines += len(line) / maxWidth
			if len(line) % maxWidth == 0 {
				wrappedLines--
			}
			longestLine = maxWidth
		}
	}

	return lineStats{longestLine, len(lines) + wrappedLines}
}

// GetText returns the current text of this DetailsView. If stripAllTags is set to true,
// any region/color tags are stripped from the text.
func (d *DetailsView) GetText(stripAllTags bool) string {
	// Note: tview.TextView.GetText appends a newline to the original text. See
	//  https://github.com/rivo/tview/blob/master/textview.go line 333 of
	//  commit 1b3174ee3d379fc32d6d5bbc63fe108a7fa8f834 and also
	//  https://github.com/rivo/tview/issues/648
	//  Additionally, the newline is not appended at the end if stripAllTags is true.
	text := d.TextView.GetText(stripAllTags)

	if !stripAllTags {
		// Note: Approach to string handling borrowed from:
		//  https://stackoverflow.com/questions/31418376/slice-unicode-ascii-strings-in-golang
		length := utf8.RuneCountInString(text)
		runes := []rune(text)
		if length > 0 && runes[length-1] == '\n' {
			return string(runes[:length-1])
		}
	}

	return text
}

// handleScrollArea calculates the width and height of the area that is scrollable in the
// DetailsView. This is a handler function that assists in drawing scroll bars on
// the DetailsView's borders.
func (d *DetailsView) handleScrollArea() (width, height int) {
	return d.LongestLine, d.LineCount
}

// handleScrollPosition calculates the current scroll position of the DirectoryList. This
// is a handler function that assists in drawing scroll bars on the DirectoryList's borders.
func (d *DetailsView) handleScrollPosition() (vScroll, hScroll int) {
	return d.GetScrollOffset()
}

// SetWordWrap sets the word wrap setting in the underlying tview.TextView and recalculates
// the content width and height.
func (d *DetailsView) SetWordWrap(wrapOnWords bool) *DetailsView {
	// Note: Characters such as '-' will wrap mid-word due to it being a line break
	//  boundary pattern, https://github.com/rivo/tview/blob/2a6de950f73bdc70658f7e754d4b5593f15c8408/util.go#L27
	d.HasWordWrap = wrapOnWords
	d.TextView.SetWordWrap(wrapOnWords)
	return d.refreshLineStats()
}

// SetText sets the content of the DetailsView and recalculates the width and height of that
// content.
func (d *DetailsView) SetText(text string) *DetailsView {
	d.TextView.SetText(text)
	return d.refreshLineStats()
}

// SetRect sets the bounds and screen location of the DetailsView.
func (d *DetailsView) SetRect(x, y, width, height int) {
	d.TextView.SetRect(x, y, width, height)
	d.refreshLineStats()
}

// Clear empties the DetailsView content and resets its title to the default value.
func (d *DetailsView) Clear() *DetailsView {
	d.TextView.Clear()
	d.SetTitle(detailsViewTitle)

	return d
}