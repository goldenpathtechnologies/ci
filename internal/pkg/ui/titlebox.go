package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func GetTitleBoxUI() *tview.TextView {
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
