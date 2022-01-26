package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/goldenpathtechnologies/ci/internal/pkg/dirctrl"
	"github.com/goldenpathtechnologies/ci/internal/pkg/options"
	"github.com/rivo/tview"
)

func setApplicationStyles() {
	// TODO: Setting this interferes with the styles of other components such
	//  as the List. Find a way to target styles to specific components.
	//  Additionally, runes display horribly in PowerShell if not using
	//  Windows Terminal. Find a way to fix this.
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
}

func Run(app *App, appOptions *options.AppOptions) error {
	setApplicationStyles()

	pages := tview.NewPages()
	filter := CreateFilterForm()
	details := CreateDetailsView()
	titleBox := CreateTitleBox()

	directoryController := dirctrl.NewDefaultDirectoryController()

	list := CreateDirectoryList(
		app,
		titleBox,
		filter,
		pages,
		details,
		directoryController,
		appOptions).
		Init()

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(titleBox, 5, 0, false).
		AddItem(tview.NewFlex().
			AddItem(list, 0, 1, true).
			AddItem(details, 0, 2, false),
			0, 1, false)

	pages.AddPage("Home", flex, true, true).
		AddPage("Filter", CreateModal(filter, 40, 7), true, false)

	if err := app.SetRoot(pages, true).SetFocus(list).Run(); err != nil {
		app.Stop()
		return err
	}

	return nil
}
