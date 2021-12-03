package ui

import (
	"ci/internal/pkg/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func SetApplicationStyles() {
	// TODO: Setting this interferes with the styles of other components such
	//  as the List. Find a way to target styles to specific components.
	//  Additionally, runes display horribly in PowerShell if not using
	//  Windows Terminal. Find a way to fix this.
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
}

func run(app *tview.Application, args []string) error {
	SetApplicationStyles()

	pages := tview.NewPages()

	filter := GetFilterUI()
	details := GetDetailsUI()
	titleBox := GetTitleBoxUI()
	list := GetListUI(app, titleBox, filter, pages, details)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(titleBox, 5, 0, false).
		AddItem(tview.NewFlex().
			AddItem(list, 0, 1, true).
			AddItem(details, 0, 2, false),
			0, 1, false)

	pages.AddPage("Home", flex, true, true).
		AddPage("Filter", CreateModalUI(filter, 40, 7), true, false)

	if err := app.SetRoot(pages, true).SetFocus(list).Run(); err != nil {
		utils.ExitScreenBuffer()

		return err
	}

	return nil
}
