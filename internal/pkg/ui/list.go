package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/goldenpathtechnologies/ci/internal/pkg/dirctrl"
	"github.com/goldenpathtechnologies/ci/internal/pkg/options"
	"github.com/rivo/tview"
	"path/filepath"
	"strings"
)

const (
	listItemQuit     = "<Quit>"
	listItemHelp     = "<Help>"
	listItemFilter   = "<Filter>"
	listItemEnterDir = "<Enter directory>"
)

const (
	listTitle        = "Directory List"
	detailsHelpTitle = "Help"
)

// DirectoryList is responsible for providing the user interface that enables users to
// quickly navigate directories and select other options.
type DirectoryList struct {
	*tview.List
	app        *App
	appOptions *options.AppOptions
	pages      *tview.Pages
	titleBox   *tview.TextView
	filter     *FilterForm
	details    *DetailsView
	dirUtil    dirctrl.DirectoryController
	currentDir string
	filterText string
	menuItems  map[string]string
}

// CreateDirectoryList creates a new instance of DirectoryList.
func CreateDirectoryList(
	app *App,
	titleBox *tview.TextView,
	filter *FilterForm,
	pages *tview.Pages,
	details *DetailsView,
	directoryController dirctrl.DirectoryController,
	appOptions *options.AppOptions,
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

// Init prepares the DirectoryList for usage by initializing data and event handlers.
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

// loadDetailsForCurrentDirectory updates the details component with the file list for the
// current active directory of the DirectoryList.
func (d *DirectoryList) loadDetailsForCurrentDirectory() {
	d.details.
		Clear().
		SetText(d.getDetailsText(d.currentDir)).
		ScrollToBeginning()
}

// getDetailsText returns the file list text that gets displayed in the Details pane.
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

// handleDetailsInputCapture is an event handler that processes key events for the details
// component of the DirectoryList.
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

// handleFilterEntry is an event handler for the DirectoryList's filter component that
// triggers when text entry is completed or cancelled.
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

// configureBorder applies default settings to the DirectoryList border and enables scroll bars.
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

// handleScrollArea calculates the width and height of the area that is scrollable in the
// DirectoryList. This is a handler function that assists in drawing scroll bars on
// the DirectoryList's borders.
func (d *DirectoryList) handleScrollArea() (width, height int) {
	_, _, listWidth, _ := d.GetInnerRect()
	listHeight := d.GetItemCount()

	return listWidth, listHeight
}

// handleScrollPosition calculates the current scroll position of the DirectoryList. This
// is a handler function that assists in drawing scroll bars on the DirectoryList's borders.
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

// configureInputEvents sets the input capture handler function for the DirectoryList.
func (d *DirectoryList) configureInputEvents() *DirectoryList {
	d.SetInputCapture(d.handleInputCapture)

	return d
}

// handleInputCapture is an event handler that processes key events for the DirectoryList.
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

// handleLeftKeyEvent handles left arrow key presses. The left arrow key navigates to the parent directory.
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

// load refreshes static menu items and the list of navigable directories.
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

// addNavigableItem adds to the DirectoryList an item that contains a directory name and selection handler.
func (d *DirectoryList) addNavigableItem(dirName string) {
	if isMatch, _ := filepath.Match(d.filterText, dirName); len(d.filterText) == 0 || isMatch {
		d.AddItem(dirName,
			"",
			0,
			d.getNavigableItemSelectionHandler(dirName))
	}
}

// getNavigableItemSelectionHandler handles the navigable item event by printing the path to the dirName
// and exiting the program in the function it returns.
func (d *DirectoryList) getNavigableItemSelectionHandler(dirName string) func() {
	return func() {
		path := d.currentDir + dirctrl.OsPathSeparator + dirName
		d.app.PrintAndExit(path)
	}
}

// handleHelpSelection handles the display of help information in the details component when the help
// list item is selected.
func (d *DirectoryList) handleHelpSelection() {
	d.setDetailsText(listItemHelp)
}

// handleRightKeyEvent handles right arrow key presses. The right arrow key navigates to the selected
// directory or indicates if the navigation is not possible due to insufficient privileges.
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

// isMenuItem determines if the supplied text equals the name of any menuItems.
func (d *DirectoryList) isMenuItem(text string) bool {
	_, exists := d.menuItems[text]
	return exists
}

// setPreviousDetailsText sets the content of the details component to the directory info of
// the previous item in the DirectoryList.
func (d *DirectoryList) setPreviousDetailsText() {
	item, _ := d.List.GetItemText(d.getNextItemIndex(false))
	d.setDetailsText(item)
}

// getNextItemIndex calculates the indices of adjacent items to the current one in the
// DirectoryList.
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

// setDetailsText sets the content of the details component depending on the dirName supplied.
// Items representing a directory will display the list of files in that directory. Menu items
// display different content depending on which one is provided to this function.
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

// setNextDetailsText sets the content of the details component to the directory info of
// the next item in the DirectoryList.
func (d *DirectoryList) setNextDetailsText() {
	item, _ := d.List.GetItemText(d.getNextItemIndex(true))
	d.setDetailsText(item)
}
