package ui

import (
	"ci/internal/pkg/utils"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/karrick/godirwalk"
	"github.com/rivo/tview"
	"path/filepath"
	"strings"
)

const (
	listUITitle = "Directory List"
	listUIQuit = "<Quit>"
	listUIHelp = "<Help>"
	listUIFilter = "<Filter>"
	listUIEnterDir = "<Enter directory>"
)

func GetListUI(
	app *App,
	titleBox *tview.TextView,
	filter *tview.InputField,
	pages *tview.Pages,
	details *tview.TextView,
) *tview.List {

	list := tview.NewList().
		ShowSecondaryText(false).
		SetSelectedTextColor(tcell.ColorBlack)

	list.SetBorder(true).
		SetTitle(listUITitle).
		SetBorderPadding(1, 1, 0, 1).
		SetDrawFunc(GetScrollBarDrawFunc(
			list,
			func() (width, height int) {
				_, _, listWidth, _ := list.GetInnerRect()
				listHeight := list.GetItemCount()

				return listWidth, listHeight
			},
			func() (vScroll, hScroll int) {
				selectedItem := list.GetCurrentItem()
				itemCount := list.GetItemCount()
				_, _, _, pageHeight := list.GetInnerRect()

				v, h := list.GetOffset()

				if selectedItem == 0 {
					v = 0
				}

				if selectedItem == itemCount-1 {
					v = itemCount - pageHeight
				}

				return v, h
			}))

	currentDir, err := utils.GetInitialDirectory()
	app.HandleError(err, true)

	titleBox.Clear()
	titleBox.SetText(currentDir)

	details.Clear().
		SetText(getDetailsText(app, currentDir)).
		ScrollToBeginning().
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyEscape:
				fallthrough
			case tcell.KeyEnter:
				fallthrough
			case tcell.KeyTab:
				app.SetFocus(list)
				return nil
			case 'q':
				app.PrintAndExit(".")
			}

			return event
		})

	var filterText string

	menuItems := map[string]string{
		listUIQuit:     listUIQuit,
		listUIHelp:     listUIHelp,
		listUIFilter:   listUIFilter,
		listUIEnterDir: listUIEnterDir,
	}

	isMenuItem := func(text string) bool {
		_, exists := menuItems[text]
		return exists
	}

	loadList := func(dir string) {
		list.Clear()

		list.AddItem(listUIEnterDir, "", 'e', func() {
			app.PrintAndExit(currentDir)
		})

		scanner, err := godirwalk.NewScanner(currentDir)
		app.HandleError(err, true)

		for scanner.Scan() {
			d, err := scanner.Dirent()
			app.HandleError(err, true)

			if d.IsDir() {
				if isMatch, _ := filepath.Match(filterText, d.Name()); len(filterText) == 0 || isMatch {
					list.AddItem(d.Name() + utils.OsPathSeparator, "", 0, func() {
						path := currentDir + d.Name() + utils.OsPathSeparator
						app.PrintAndExit(path)
					})
				}
			}
		}

		list.AddItem(listUIQuit, "Press to exit", 'q', func() {
			app.PrintAndExit(".")
		})

		list.AddItem(listUIHelp, "Get help with this program", 'h', func(){})

		list.AddItem(listUIFilter, "Filter directories by text", 'f', func() {
			pages.ShowPage("Filter")
			app.SetFocus(filter)
		})

		titleBox.Clear()
		titleBox.SetText(currentDir)
	}
	loadList(currentDir)

	// TODO: Ensure Esc does not apply any existing filter
	filter.SetDoneFunc(func(key tcell.Key) {
		filterText = filter.GetText()
		if len(filterText) > 0 {
			list.SetTitle(fmt.Sprintf("%v - Filter: %v", listUITitle, filterText))
		} else {
			list.SetTitle(listUITitle)
		}
		filter.SetText("")
		pages.HidePage("Filter")
		app.SetFocus(list)
		loadList(currentDir)
	})

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		selectedItem, _ := list.GetItemText(list.GetCurrentItem())

		getNextItemIndex := func(isIncrementing bool) int {
			var increment int
			if isIncrementing {
				increment = 1
			} else {
				increment = -1
			}
			itemCount := list.GetItemCount()
			nextItemIndex := list.GetCurrentItem() + increment
			// Note: Euclidean modulo operation, https://stackoverflow.com/questions/43018206/modulo-of-negative-integers-in-go
			return ((nextItemIndex % itemCount) + itemCount) % itemCount
		}

		displayDirectoryDetails := func(isNavigatingDown bool) {
			details.Clear()
			nextItemIndex := getNextItemIndex(isNavigatingDown)
			nextItem, _ := list.GetItemText(nextItemIndex)
			if !isMenuItem(nextItem) {
				details.SetText(getDetailsText(app, currentDir + nextItem))
			} else if nextItem == listUIEnterDir {
				details.SetText(getDetailsText(app, currentDir))
			}
			details.ScrollToBeginning()
		}

		switch event.Key() {
		case tcell.KeyLeft:
			if strings.Count(currentDir, utils.OsPathSeparator) > 1 {
				filterText = ""
				list.SetTitle(listUITitle)
				paths := strings.Split(currentDir, utils.OsPathSeparator)
				paths = paths[:len(paths)-2]
				currentDir = strings.Join(paths, utils.OsPathSeparator) + utils.OsPathSeparator
				loadList(currentDir)

				details.Clear()
				details.SetText(getDetailsText(app, currentDir)).
					ScrollToBeginning()
			}
			return nil
		case tcell.KeyRight:
			if !isMenuItem(selectedItem) {
				filterText = ""
				list.SetTitle(listUITitle)
				nextDir := currentDir + selectedItem
				if utils.DirectoryIsAccessible(nextDir) {
					currentDir = nextDir
					loadList(currentDir)
				} else {
					details.Clear().
						SetText("[red]Directory inaccessible, unable to navigate. You may have insufficient privileges.[white]").
						ScrollToBeginning()
				}
			}
			return nil
		case tcell.KeyUp:
			displayDirectoryDetails(false)
			return event
		case tcell.KeyDown:
			displayDirectoryDetails(true)
			return event
		case tcell.KeyTab:
			app.SetFocus(details)
			return nil
		}

		return event
	})

	return list
}

func getDetailsText(app *App, directory string) string {
	var (
		detailsText string
		err         error
	)

	if detailsText, err = utils.GetDirectoryInfo(directory); err != nil {
		dErr, isDirError := err.(*utils.DirectoryError)
		if isDirError && dErr.ErrorCode == utils.DirUnprivilegedError {
			detailsText = "[red]Unable to read directory details. You may have insufficient privileges.[white]"
		} else {
			app.HandleError(err, true)
		}
	}

	return detailsText
}
