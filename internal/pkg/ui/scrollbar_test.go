package ui

import (
	td "github.com/goldenpathtechnologies/ci/testdata"
	"github.com/goldenpathtechnologies/ci/testdata/utils"
	"github.com/rivo/tview"
	"strings"
	"testing"
)

const (
	screenWidth = 80
	screenHeight = 50
)

type padding struct {
	top, bottom, left, right int
}

func setUpTestTextView(
	scrollableAreaText string,
	width,
	height int,
	border bool,
	wrap bool,
	padding padding,
) *tview.TextView {

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(wrap).
		SetText(scrollableAreaText)
	textView.
		SetBorder(border).
		SetBorderPadding(padding.top, padding.bottom, padding.left, padding.right).
		SetRect(0, 0, width, height)

	return textView
}

func Test_ScrollBar_GetScrollBarDrawFunc_ScrollbarNotDrawnWithoutBorderOrPadding(t *testing.T) {
	runScrollBarBorderAndPaddingTests(
		false,
		padding{0, 0, 0, 0},
		func(output, scrollBarChars string) {
			if strings.ContainsAny(output, scrollBarChars) {
				t.Errorf("Expected the scroll bar not to be drawn. Output:\n%s\n", output)
			}
		})
}

func runScrollBarBorderAndPaddingTests(border bool, padding padding, assert func(output, scrollBarChars string)) {
	testApp := utils.NewTestApp().Init(screenWidth, screenHeight)
	testData := td.TestText["LoremIpsum"]
	view := setUpTestTextView(
		testData.Text,
		testData.LongestWrappedLine,
		testData.LineCount-1,
		border,
		false,
		padding)

	scrollAreaFunc := func() (width, height int) {
		return testData.LongestLine, testData.LineCount
	}

	scrollPosFunc := func() (vScroll, hScroll int) {
		return view.GetScrollOffset()
	}

	view.
		SetDrawFunc(
			GetScrollBarDrawFunc(
				view,
				scrollAreaFunc,
				scrollPosFunc))

	testApp.Run(view, true, func() {
		output := testApp.GetPrimitiveOutput()
		scrollChars :=
			string(tview.Borders.Vertical) +
				string(tview.Borders.Horizontal) +
				string(tview.Borders.VerticalFocus) +
				string(tview.Borders.HorizontalFocus)

		assert(output, scrollChars)
	})
}

func Test_ScrollBar_GetScrollBarDrawFunc_ScrollbarDrawnWithPadding(t *testing.T) {
	runScrollBarBorderAndPaddingTests(
		false,
		padding{1, 1, 1, 1},
		func(output, scrollBarChars string) {
			if !strings.ContainsAny(output, scrollBarChars) {
				t.Errorf("Expected scroll bar to be drawn. Output:\n%s\n", output)
			}
		})
}

func Test_ScrollBar_GetScrollBarDrawFunc_ScrollbarDrawnWithBorder(t *testing.T) {
	runScrollBarBorderAndPaddingTests(
		true,
		padding{0, 0, 0, 0},
		func(output, scrollBarChars string) {
			if !strings.ContainsAny(output, scrollBarChars) {
				t.Errorf("Expected scroll bar to be drawn. Output:\n%s\n", output)
			}
		})
}
