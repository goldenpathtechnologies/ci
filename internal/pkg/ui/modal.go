package ui

import "github.com/rivo/tview"

func CreateModalUI(widget tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(widget, 1, 1, 1, 1, 0, 0, false)
}
