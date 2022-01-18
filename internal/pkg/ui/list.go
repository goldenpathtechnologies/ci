package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/goldenpathtechnologies/ci/internal/pkg/dirctrl"
	"github.com/goldenpathtechnologies/ci/internal/pkg/flags"
	"github.com/rivo/tview"
	"path/filepath"
	"strings"
)

const (
	listTitle           = "Directory List"
	listItemQuit        = "<Quit>"
	listItemHelp        = "<Help>"
	listItemFilter      = "<Filter>"
	listItemEnterDir    = "<Enter directory>"
	detailsHelpTitle    = "Help"
)

type DirectoryList struct {
	*tview.List
	app        *App
	appOptions *flags.AppOptions
	pages      *tview.Pages
	titleBox   *tview.TextView
	filter     *FilterForm
	details    *DetailsView
	dirUtil    dirctrl.DirectoryController
	currentDir string
	filterText string
	menuItems  map[string]string
}

func CreateDirectoryList(
	app *App,
	titleBox *tview.TextView,
	filter *FilterForm,
	pages *tview.Pages,
	details *DetailsView,
	directoryController dirctrl.DirectoryController,
	appOptions *flags.AppOptions,
) *DirectoryList {
	list := tview.NewList().
		ShowSecondaryText(false).
		SetSelectedTextColor(tcell.ColorBlack)

	menuItems := map[string]string{
		listItemQuit:     listItemQuit,
		listItemHelp:     listItemHelp,
		listItemFilter:   listItemFilter,
		listItemEnterDir: listItemEnterDir,
	}

	return &DirectoryList{
		List:       list,
		app:        app,
		appOptions: appOptions,
		pages:      pages,
		titleBox:   titleBox,
		filter:     filter,
		details:    details,
		dirUtil:    directoryController,
		menuItems:  menuItems,
	}
}

func (d *DirectoryList) Init() *DirectoryList {
	var err error

	d.currentDir, err = d.dirUtil.GetInitialDirectory()
	d.app.HandleError(err, true)

	d.titleBox.Clear()
	d.titleBox.SetText(d.currentDir)

	d.loadDetailsForCurrentDirectory()
	d.details.SetInputCapture(d.handleDetailsInputCapture)

	d.filter.SetDoneHandler(d.handleFilterEntry)

	d.configureBorder().configureInputEvents().load()

	return d
}

func (d *DirectoryList) loadDetailsForCurrentDirectory() {
	d.details.
		Clear().
		SetText(d.getDetailsText(d.currentDir)).
		ScrollToBeginning()
}

func (d *DirectoryList) getDetailsText(directory string) string {
	var (
		detailsText string
		err         error
	)

	if detailsText, err = d.dirUtil.GetDirectoryInfo(directory); err != nil {
		dErr, isDirError := err.(*dirctrl.DirectoryError)
		if isDirError && dErr.ErrorCode == dirctrl.DirUnprivilegedError {
			detailsText = "[red]Unable to read directory details. You may have insufficient privileges.[white]"
		} else {
			d.app.HandleError(err, true)
		}
	}

	return detailsText
}

func (d *DirectoryList) handleDetailsInputCapture(event *tcell.EventKey) *tcell.EventKey {
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
		return nil
	}

	return event
}

func (d *DirectoryList) handleFilterEntry(key tcell.Key) {
	if key == tcell.KeyEsc {
		d.filter.Clear()
	}

	d.filterText = d.filter.GetText()

	if len(d.filterText) > 0 {
		d.SetTitle(fmt.Sprintf("%v - Filter: %v", listTitle, d.filterText))
	} else {
		d.SetTitle(listTitle)
	}

	d.filter.Clear()
	d.pages.HidePage("Filter")
	d.app.SetFocus(d)
	d.load()
}

func (d *DirectoryList) configureBorder() *DirectoryList {
	d.SetBorder(true).
		SetTitle(listTitle).
		SetBorderPadding(1, 1, 0, 1).
		SetDrawFunc(GetScrollBarDrawFunc(
			d,
			d.handleScrollArea,
			d.handleScrollPosition))

	return d
}

func (d *DirectoryList) handleScrollArea() (width, height int) {
	_, _, listWidth, _ := d.GetInnerRect()
	listHeight := d.GetItemCount()

	return listWidth, listHeight
}

func (d *DirectoryList) handleScrollPosition() (vScroll, hScroll int) {
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
}

func (d *DirectoryList) configureInputEvents() *DirectoryList {
	d.SetInputCapture(d.handleInputCapture)

	return d
}

func (d *DirectoryList) handleInputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyLeft:
		d.handleLeftKeyEvent()
		return nil
	case tcell.KeyRight:
		d.handleRightKeyEvent()
		return nil
	case tcell.KeyUp:
		d.setPreviousDetailsText()
		return event
	case tcell.KeyDown:
		d.setNextDetailsText()
		return event
	case tcell.KeyTab:
		d.app.SetFocus(d.details)
		return nil
	}

	return event
}

func (d *DirectoryList) handleLeftKeyEvent() {
	paths := strings.Split(strings.TrimRight(d.currentDir, dirctrl.OsPathSeparator), dirctrl.OsPathSeparator)
	if len(paths) > 1 {
		d.filterText = ""
		d.SetTitle(listTitle)
		paths = paths[:len(paths)-1]
		if len(paths) == 1 && (paths[0] == "" || strings.Contains(paths[0], ":")) {
			d.currentDir, _ = d.dirUtil.GetAbsolutePath(dirctrl.OsPathSeparator)
		} else {
			d.currentDir = strings.Join(paths, dirctrl.OsPathSeparator)
		}
		d.load()
		d.loadDetailsForCurrentDirectory()
	}
}

func (d *DirectoryList) load() {
	d.Clear()

	d.AddItem(listItemEnterDir, "", 'e', func() {
		d.app.PrintAndExit(d.currentDir)
	})

	if err := d.dirUtil.ScanDirectory(d.currentDir, func(dirName string) {
		d.addNavigableItem(dirName)
	}); err != nil {
		d.app.HandleError(err, true)
	}

	d.AddItem(listItemFilter, "Filter directories by text", 'f', func() {
		d.pages.ShowPage("Filter")
		d.app.SetFocus(d.filter)
	})

	d.AddItem(
		listItemHelp,
		"Get help with this program",
		'h',
		d.handleHelpSelection)

	d.AddItem(listItemQuit, "Press to exit", 'q', func() {
		d.app.PrintAndExit(".")
	})

	d.titleBox.Clear()
	d.titleBox.SetText(d.currentDir)
}

func (d *DirectoryList) addNavigableItem(dirName string) {
	if isMatch, _ := filepath.Match(d.filterText, dirName); len(d.filterText) == 0 || isMatch {
		d.AddItem(dirName,
			"",
			0,
			d.getNavigableItemSelectionHandler(dirName))
	}
}

func (d *DirectoryList) getNavigableItemSelectionHandler(dirName string) func() {
	return func() {
		path := d.currentDir + dirctrl.OsPathSeparator + dirName
		d.app.PrintAndExit(path)
	}
}

func (d *DirectoryList) handleHelpSelection() {
	d.setDetailsText(listItemHelp)
}

func (d *DirectoryList) handleRightKeyEvent() {
	selectedItem, _ := d.GetItemText(d.GetCurrentItem())

	if !d.isMenuItem(selectedItem) {
		d.filterText = ""
		d.SetTitle(listTitle)
		pathCount := len(strings.Split(strings.TrimRight(d.currentDir, dirctrl.OsPathSeparator), dirctrl.OsPathSeparator))
		var pathSeparator string
		if pathCount > 1 {
			pathSeparator = dirctrl.OsPathSeparator
		} else {
			pathSeparator = ""
		}
		nextDir := d.currentDir + pathSeparator + selectedItem
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

func (d *DirectoryList) setPreviousDetailsText() {
	item, _ := d.List.GetItemText(d.getNextItemIndex(false))
	d.setDetailsText(item)
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

func (d *DirectoryList) setDetailsText(dirName string) {
	d.details.Clear()
	if !d.isMenuItem(dirName) {
		d.details.SetText(d.getDetailsText(d.currentDir + dirctrl.OsPathSeparator + dirName))
	} else if dirName == listItemEnterDir {
		d.details.SetText(d.getDetailsText(d.currentDir))
	} else if dirName == listItemHelp {
		d.details.SetText(GetHelpText(d.appOptions))
		d.details.SetTitle(detailsHelpTitle)
	}
	d.details.ScrollToBeginning()
}

func (d *DirectoryList) setNextDetailsText() {
	item, _ := d.List.GetItemText(d.getNextItemIndex(true))
	d.setDetailsText(item)
}
