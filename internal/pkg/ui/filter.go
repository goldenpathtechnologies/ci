package ui

import "github.com/rivo/tview"

const maxFilterLength = 32

func CreateFilterPane() *tview.InputField {
	filter := tview.NewInputField().
		SetLabel("Enter filter text: ").
		SetAcceptanceFunc(handleFilterAcceptance)

	filter.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetTitle("Filter Directory List")

	return filter
}

func handleFilterAcceptance(textToCheck string, lastChar rune) bool {
	if lastChar == '/' || lastChar == '\\' {
		return false
	}

	return len(textToCheck) <= maxFilterLength
}
