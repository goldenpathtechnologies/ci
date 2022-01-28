package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CreateTitleBox creates and configures the title box of the application that displays
// the current navigated directory.
func CreateTitleBox() *tview.TextView {
	titleBox := tview.NewTextView().
		SetText("No directory has been selected yet!").
		SetTextAlign(tview.AlignCenter).
		SetScrollable(false)

	titleBox.SetBorder(true).
		SetTitle(`ci - Interactive cd`).
		SetBorderPadding(1,1,1,1).
		SetTitleColor(tcell.ColorGreen)

	return titleBox
}
