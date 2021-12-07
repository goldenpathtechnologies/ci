package ui

import "github.com/rivo/tview"

func CreateFilterPane() *tview.InputField {
	// TODO: Put the field width and max length in a constant
	filter := tview.NewInputField().
		SetLabel("Enter filter text: ").
		SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
			if lastChar == '/' || lastChar == '\\' {
				return false
			}

			return len(textToCheck) <= 32
		})

	filter.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetTitle("Filter Directory List")

	return filter
}

