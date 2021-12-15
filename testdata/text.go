package testdata

type TextData struct {
	Text                   string
	LongestLine            int
	LineCount              int
	LongestWrappedLine     int
	LongestWordWrappedLine int
	WrappedLineCount       int
	WordWrappedLineCount   int
	ViewWidth              int
}

var (
	TestText = map[string]TextData{
		"Empty": {
			Text: "",
			LongestLine: 0,
			LineCount: 1,
			LongestWrappedLine: 0,
			LongestWordWrappedLine: 0,
			WrappedLineCount: 1,
			WordWrappedLineCount: 1,
			ViewWidth: 10,
		},
		"LoremIpsum": {
			Text: `Lorem ipsum dolor sit amet, consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in
voluptate velit esse cillum dolore eu fugiat nulla pariatur.
Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia
deserunt mollit anim id est laborum.`,
			LongestLine:            74,
			LineCount:              7,
			LongestWrappedLine:     64,
			LongestWordWrappedLine: 61,
			WrappedLineCount:       11,
			WordWrappedLineCount:   11,
			ViewWidth:              64,
		},
		"TrailingNewline": {
			Text: `This is some test text
that has some newlines and
also one at the end
`,
			LongestLine: 26,
			LineCount: 4,
			LongestWrappedLine: 26,
			LongestWordWrappedLine: 26,
			WrappedLineCount: 4,
			WordWrappedLineCount: 4,
			ViewWidth: 50,
		},
		"TrailingSpace": {
			Text: "This is some text that just has a space at the end ",
			LongestLine: 51,
			LineCount: 1,
			LongestWrappedLine: 51,
			LongestWordWrappedLine: 51,
			WrappedLineCount: 1,
			WordWrappedLineCount: 1,
			ViewWidth: 60,
		},
		"SingleLineOnceWrapped": {
			Text: "bbbbbbbbbbbbbbbb",
			LongestLine: 16,
			LineCount: 1,
			LongestWrappedLine: 8,
			LongestWordWrappedLine: 8,
			WrappedLineCount: 2,
			WordWrappedLineCount: 2,
			ViewWidth: 8,
		},
		"SingleLineThriceWrapped": {
			Text: "hhhhhhhhhhhhhhhhhhhhhhhh",
			LongestLine: 24,
			LineCount: 1,
			LongestWrappedLine: 6,
			LongestWordWrappedLine: 6,
			WrappedLineCount: 4,
			WordWrappedLineCount: 4,
			ViewWidth: 6,
		},
		"FourEqualLinesEachWrapped": {
			Text: "aaaa\naaaa\naaaa\naaaa",
			LongestLine: 4,
			LineCount: 4,
			LongestWrappedLine: 2,
			LongestWordWrappedLine: 2,
			WrappedLineCount: 8,
			WordWrappedLineCount: 8,
			ViewWidth: 2,
		},
		"ThreeUnequalLinesEachWrapped": {
			Text: "aaaaa\naaaa\naaaa",
			LongestLine: 5,
			LineCount: 3,
			LongestWrappedLine: 2,
			LongestWordWrappedLine: 2,
			WrappedLineCount: 7,
			WordWrappedLineCount: 7,
			ViewWidth: 2,
		},
		// TODO: Add test data with tabs at the end. Currently, tview converts tabs to spaces. See,
		//  https://github.com/rivo/tview/blob/2a6de950f73bdc70658f7e754d4b5593f15c8408/textview.go#L663
	}
)

func RunTextTestCases(testFunc func(data TextData, name string)) {
	for i, t := range TestText {
		testFunc(t, i)
	}
}