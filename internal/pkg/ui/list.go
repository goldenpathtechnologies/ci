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
	details    *DetailsView
	dirUtil    utils.DirectoryController
	currentDir string
	filterText string
	menuItems  map[string]string
}

func CreateDirectoryList(
	app *App,
	titleBox *tview.TextView,
	filter *tview.InputField,
	pages *tview.Pages,
	details *DetailsView,
) *DirectoryList {

	list, err := newDirectoryList(app, titleBox, filter, pages, details)
	app.HandleError(err, true)

	list.titleBox.Clear()
	list.titleBox.SetText(list.currentDir)

	list.details.Clear().
		SetText(list.getDetailsText(list.currentDir)).
		ScrollToBeginning().
		SetInputCapture(list.getDetailsInputCaptureHandler())

	list.filter.SetDoneFunc(list.getFilterEntryHandler())

	list.
		configureBorder().
		configureInputEvents().
		load()

	return list
}


func newDirectoryList(
	app *App,
	titleBox *tview.TextView,
	filter *tview.InputField,
	pages *tview.Pages,
	details *DetailsView,
) (*DirectoryList, error) {
	var (
		currentDir string
		err        error
	)

	list := tview.NewList().
		ShowSecondaryText(false).
		SetSelectedTextColor(tcell.ColorBlack)

	dirUtil := utils.NewDefaultDirectoryController()

	currentDir, err = dirUtil.GetInitialDirectory()

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
		details:    details,
		dirUtil:    dirUtil,
		currentDir: currentDir,
		menuItems:  menuItems,
	}, err
}

func (d *DirectoryList) getDetailsText(directory string) string {
	var (
		detailsText string
		err         error
	)

	if detailsText, err = d.dirUtil.GetDirectoryInfo(directory); err != nil {
		dErr, isDirError := err.(*utils.DirectoryError)
		if isDirError && dErr.ErrorCode == utils.DirUnprivilegedError {
			detailsText = "[red]Unable to read directory details. You may have insufficient privileges.[white]"
		} else {
			d.app.HandleError(err, true)
		}
	}

	return detailsText
}

func (d *DirectoryList) getDetailsInputCaptureHandler() func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			fallthrough
		case tcell.KeyEnter:
			fallthrough
		case tcell.KeyTab:
			d.app.SetFocus(d)
			return nil
		case 'q':
			d.app.PrintAndExit(".")
		}

		return event
	}
}

func (d *DirectoryList) getFilterEntryHandler() func(key tcell.Key) {
	return func(key tcell.Key) {
		// TODO: Ensure Esc does not apply any existing filter
		d.filterText = d.filter.GetText()
		if len(d.filterText) > 0 {
			d.SetTitle(fmt.Sprintf("%v - Filter: %v", listUITitle, d.filterText))
		} else {
			d.SetTitle(listUITitle)
		}
		d.filter.SetText("")
		d.pages.HidePage("Filter")
		d.app.SetFocus(d)
		d.load()
	}
}

func (d *DirectoryList) configureBorder() *DirectoryList {
	d.SetBorder(true).
		SetTitle(listUITitle).
		SetBorderPadding(1, 1, 0, 1).
		SetDrawFunc(GetScrollBarDrawFunc(
			d,
			getListScrollAreaHandler(d),
			getListScrollPositionHandler(d)))

	return d
}

func getListScrollAreaHandler(list *DirectoryList) func() (width, height int) {
	return func() (width, height int) {
			_, _, listWidth, _ := list.GetInnerRect()
			listHeight := list.GetItemCount()

			return listWidth, listHeight
	}
}

func getListScrollPositionHandler(list *DirectoryList) func() (vScroll, hScroll int) {
	return func() (vScroll, hScroll int) {
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
	}
}

func (d *DirectoryList) configureInputEvents() *DirectoryList {
	d.SetInputCapture(d.getInputCaptureHandler())

	return d
}

func (d *DirectoryList) getInputCaptureHandler() func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			d.handleLeftKeyEvent()
			return nil
		case tcell.KeyRight:
			d.handleRightKeyEvent()
			return nil
		case tcell.KeyUp:
			d.setDetailsText(false)
			return event
		case tcell.KeyDown:
			d.setDetailsText(true)
			return event
		case tcell.KeyTab:
			d.app.SetFocus(d.details)
			return nil
		}

		return event
	}
}

func (d *DirectoryList) handleLeftKeyEvent() {
	if strings.Count(d.currentDir, utils.OsPathSeparator) > 1 {
		d.filterText = ""
		d.SetTitle(listUITitle)
		paths := strings.Split(d.currentDir, utils.OsPathSeparator)
		paths = paths[:len(paths)-2]
		d.currentDir = strings.Join(paths, utils.OsPathSeparator) + utils.OsPathSeparator
		d.load()

		d.details.Clear()
		d.details.
			SetText(d.getDetailsText(d.currentDir)).
			ScrollToBeginning()
	}
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

func (d *DirectoryList) handleRightKeyEvent() {
	selectedItem, _ := d.GetItemText(d.GetCurrentItem())

	if !d.isMenuItem(selectedItem) {
		d.filterText = ""
		d.SetTitle(listUITitle)
		nextDir := d.currentDir + selectedItem
		if d.dirUtil.DirectoryIsAccessible(nextDir) {
			d.currentDir = nextDir
			d.load()
		} else {
			d.details.Clear()
			d.details.SetText("[red]Directory inaccessible, unable to navigate. You may have insufficient privileges.[white]").
				ScrollToBeginning()
		}
	}
}

func (d *DirectoryList) isMenuItem(text string) bool {
	_, exists := d.menuItems[text]
	return exists
}

func (d *DirectoryList) setDetailsText(isNavigatingDown bool) {
	d.details.Clear()
	nextItemIndex := d.getNextItemIndex(isNavigatingDown)
	nextItem, _ := d.GetItemText(nextItemIndex)
	if !d.isMenuItem(nextItem) {
		d.details.SetText(d.getDetailsText(d.currentDir + nextItem))
	} else if nextItem == listUIEnterDir {
		d.details.SetText(d.getDetailsText(d.currentDir))
	}
	d.details.ScrollToBeginning()
}

func (d *DirectoryList) getNextItemIndex(isIncrementing bool) int {
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
