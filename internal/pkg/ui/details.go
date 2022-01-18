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

type lineStats struct {
	longest int
	count   int
}

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

func (d *DetailsView) SetWrap(wrap bool) *DetailsView {
	d.HasWrap = wrap
	d.TextView.SetWrap(wrap)
	return d.refreshLineStats()
}

func (d *DetailsView) refreshLineStats() *DetailsView {
	_, _, viewWidth, _ := d.TextView.GetInnerRect()
	lineData := calculateLineStats(d.GetText(true), viewWidth, d.HasWrap, d.HasWordWrap)
	d.LongestLine = lineData.longest
	d.LineCount = lineData.count
	return d
}

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

func (d *DetailsView) handleScrollArea() (width, height int) {
	return d.LongestLine, d.LineCount
}

func (d *DetailsView) handleScrollPosition() (vScroll, hScroll int) {
	return d.GetScrollOffset()
}

func (d *DetailsView) SetWordWrap(wrapOnWords bool) *DetailsView {
	// Note: Characters such as '-' will wrap mid-word due to it being a line break
	//  boundary pattern, https://github.com/rivo/tview/blob/2a6de950f73bdc70658f7e754d4b5593f15c8408/util.go#L27
	d.HasWordWrap = wrapOnWords
	d.TextView.SetWordWrap(wrapOnWords)
	return d.refreshLineStats()
}

func (d *DetailsView) SetText(text string) *DetailsView {
	d.TextView.SetText(text)
	return d.refreshLineStats()
}

func (d *DetailsView) SetRect(x, y, width, height int) {
	d.TextView.SetRect(x, y, width, height)
	d.refreshLineStats()
}

func (d *DetailsView) Clear() *DetailsView {
	d.TextView.Clear()
	d.SetTitle(detailsViewTitle)

	return d
}