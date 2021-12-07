package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/goldenpathtechnologies/ci/internal/pkg/utils"
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

type DirectoryList struct {
	*tview.List
	app        *App
	pages      *tview.Pages
	titleBox   *tview.TextView
	filter     *tview.InputField
	currentDir string
	filterText string
	menuItems  map[string]string
}

func newDirectoryList(
	app *App,
	titleBox *tview.TextView,
	filter *tview.InputField,
	pages *tview.Pages,
) (*DirectoryList, error) {
	var (
		currentDir string
		err        error
	)

	list := tview.NewList().
		ShowSecondaryText(false).
		SetSelectedTextColor(tcell.ColorBlack)

	currentDir, err = utils.GetInitialDirectory()

	menuItems := map[string]string{
		listUIQuit:     listUIQuit,
		listUIHelp:     listUIHelp,
		listUIFilter:   listUIFilter,
		listUIEnterDir: listUIEnterDir,
	}

	return &DirectoryList{
		List:       list,
		app:        app,
		pages:      pages,
		titleBox:   titleBox,
		filter:     filter,
		currentDir: currentDir,
		menuItems:  menuItems,
	}, err
}

func (d *DirectoryList) configureBorder() {
	d.SetBorder(true).
		SetTitle(listUITitle).
		SetBorderPadding(1, 1, 0, 1).
		SetDrawFunc(GetScrollBarDrawFunc(
			d,
			func() (width, height int) {
				_, _, listWidth, _ := d.GetInnerRect()
				listHeight := d.GetItemCount()

				return listWidth, listHeight
			},
			func() (vScroll, hScroll int) {
				selectedItem := d.GetCurrentItem()
				itemCount := d.GetItemCount()
				_, _, _, pageHeight := d.GetInnerRect()

				v, h := d.GetOffset()

				if selectedItem == 0 {
					v = 0
				}

				if selectedItem == itemCount-1 {
					v = itemCount - pageHeight
				}

				return v, h
			}))
}

func (d *DirectoryList) isMenuItem(text string) bool {
	_, exists := d.menuItems[text]
	return exists
}

func (d *DirectoryList) load() {
	d.Clear()

	d.AddItem(listUIEnterDir, "", 'e', func() {
		d.app.PrintAndExit(d.currentDir)
	})

	scanner, err := godirwalk.NewScanner(d.currentDir)
	d.app.HandleError(err, true)

	for scanner.Scan() {
		entry, err := scanner.Dirent()
		d.app.HandleError(err, true)

		if entry.IsDir() {
			if isMatch, _ := filepath.Match(d.filterText, entry.Name()); len(d.filterText) == 0 || isMatch {
				d.AddItem(entry.Name() + utils.OsPathSeparator, "", 0, func() {
					path := d.currentDir + entry.Name() + utils.OsPathSeparator
					d.app.PrintAndExit(path)
				})
			}
		}
	}

	d.AddItem(listUIQuit, "Press to exit", 'q', func() {
		d.app.PrintAndExit(".")
	})

	// TODO: Implement in-app help.
	d.AddItem(listUIHelp, "Get help with this program", 'h', func(){})

	d.AddItem(listUIFilter, "Filter directories by text", 'f', func() {
		d.pages.ShowPage("Filter")
		d.app.SetFocus(d.filter)
	})

	d.titleBox.Clear()
	d.titleBox.SetText(d.currentDir)
}

func (d *DirectoryList) configureInputEvents(
	app *App,
	details *tview.TextView,
) {
	d.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		selectedItem, _ := d.GetItemText(d.GetCurrentItem())

		getNextItemIndex := func(isIncrementing bool) int {
			var increment int
			if isIncrementing {
				increment = 1
			} else {
				increment = -1
			}
			itemCount := d.GetItemCount()
			nextItemIndex := d.GetCurrentItem() + increment
			// Note: Euclidean modulo operation, https://stackoverflow.com/questions/43018206/modulo-of-negative-integers-in-go
			return ((nextItemIndex % itemCount) + itemCount) % itemCount
		}

		displayDirectoryDetails := func(isNavigatingDown bool) {
			details.Clear()
			nextItemIndex := getNextItemIndex(isNavigatingDown)
			nextItem, _ := d.GetItemText(nextItemIndex)
			if !d.isMenuItem(nextItem) {
				details.SetText(getDetailsText(app, d.currentDir + nextItem))
			} else if nextItem == listUIEnterDir {
				details.SetText(getDetailsText(app, d.currentDir))
			}
			details.ScrollToBeginning()
		}

		switch event.Key() {
		case tcell.KeyLeft:
			if strings.Count(d.currentDir, utils.OsPathSeparator) > 1 {
				d.filterText = ""
				d.SetTitle(listUITitle)
				paths := strings.Split(d.currentDir, utils.OsPathSeparator)
				paths = paths[:len(paths)-2]
				d.currentDir = strings.Join(paths, utils.OsPathSeparator) + utils.OsPathSeparator
				d.load()

				details.Clear()
				details.SetText(getDetailsText(app, d.currentDir)).
					ScrollToBeginning()
			}
			return nil
		case tcell.KeyRight:
			if !d.isMenuItem(selectedItem) {
				d.filterText = ""
				d.SetTitle(listUITitle)
				nextDir := d.currentDir + selectedItem
				if utils.DirectoryIsAccessible(nextDir) {
					d.currentDir = nextDir
					d.load()
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

func CreateDirectoryList(
	app *App,
	titleBox *tview.TextView,
	filter *tview.InputField,
	pages *tview.Pages,
	details *tview.TextView,
) *DirectoryList {

	list, err := newDirectoryList(app, titleBox, filter, pages)
	app.HandleError(err, true)

	titleBox.Clear()
	titleBox.SetText(list.currentDir)

	details.Clear().
		SetText(getDetailsText(app, list.currentDir)).
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

	// TODO: Ensure Esc does not apply any existing filter
	filter.SetDoneFunc(func(key tcell.Key) {
		list.filterText = filter.GetText()
		if len(list.filterText) > 0 {
			list.SetTitle(fmt.Sprintf("%v - Filter: %v", listUITitle, list.filterText))
		} else {
			list.SetTitle(listUITitle)
		}
		filter.SetText("")
		pages.HidePage("Filter")
		app.SetFocus(list)
		list.load()
	})

	list.configureBorder()
	list.configureInputEvents(app, details)
	list.load()

	return list
}
