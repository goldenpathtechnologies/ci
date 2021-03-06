package ui

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/goldenpathtechnologies/ci/internal/pkg/dirctrl"
	"github.com/goldenpathtechnologies/ci/testdata/mock"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func Test_DirectoryList_getDetailsText_ReturnsDirectoryDetails(t *testing.T) {
	var files []fs.FileInfo

	for i := 0; i < 10; i++ {
		files = append(files, mock.GenerateMockFile())
	}

	dirCtrl := dirctrl.NewDefaultDirectoryController()
	dirCtrl.Commands = &mock.DirectoryCommands{
		ReadDirectoryFunc: func(dirname string) ([]fs.FileInfo, error) {
			return files, nil
		},
		GetAbsolutePathFunc: func(path string) (string, error) {
			return "", nil
		},
	}

	list := CreateDirectoryList(nil, nil, nil, nil, nil, dirCtrl, nil)

	detailsText := list.getDetailsText(".")

	for _, file := range files {
		if !strings.Contains(detailsText, file.Name()) {
			t.Errorf(
				"Expected the file '%s' to be in output", file.Name())
		}
	}

	if t.Failed() {
		fmt.Printf("Output:\n%s\n", detailsText)
	}
}

func Test_DirectoryList_getDetailsText_ReturnsUnprivilegedMessageWhenDirectoryInaccessible(t *testing.T) {
	dirCtrl := dirctrl.NewDefaultDirectoryController()
	dirCtrl.Commands = &mock.DirectoryCommands{
		ReadDirectoryFunc: func(dirname string) ([]fs.FileInfo, error) {
			return nil, errors.New("unable to access directory")
		},
		GetAbsolutePathFunc: func(path string) (string, error) {
			return "", nil
		},
	}

	list := CreateDirectoryList(nil, nil, nil, nil, nil, dirCtrl, nil)

	result := list.getDetailsText(".")
	expected := "[red]Unable to read directory details. You may have insufficient privileges.[white]"

	if result != expected {
		t.Errorf("Expected output to be the following:\n%s\n\nGot the following instead:\n%s\n",
			expected, result)
	}
}

func Test_DirectoryList_getDetailsText_HandlesUnexpectedErrors(t *testing.T) {
	var (
		files []fs.FileInfo
		out   bytes.Buffer
	)

	errorMessage := "error triggered by test"

	for i := 0; i < 10; i++ {
		files = append(files, mock.GenerateMockFile())
	}

	dirCtrl := dirctrl.NewDefaultDirectoryController()
	dirCtrl.Commands = &mock.DirectoryCommands{
		ReadDirectoryFunc: func(dirname string) ([]fs.FileInfo, error) {
			return files, nil
		},
		GetAbsolutePathFunc: func(path string) (string, error) {
			return "", nil
		},
	}
	dirCtrl.Writer = &mock.InfoWriter{
		WriteFunc: func(p []byte) (n int, err error) {
			return 0, errors.New(errorMessage)
		},
		FlushFunc: func() (string, error) {
			return "", nil
		},
	}

	screen := tcell.NewSimulationScreen("") // "" = UTF-8 charset
	app := NewApp(screen, io.Discard, io.Discard)
	app.handleErrorExit = func() {
		// Do nothing for test
	}

	list := CreateDirectoryList(app, nil, nil, nil, nil, dirCtrl, nil)

	// TODO: When the App struct implements a logging flag, get rid of this statement
	//  and expect error output from the default errorStream instead. Currently, the
	//  errorStream prints to io.Discard.
	log.SetOutput(&out)
	list.getDetailsText(".")

	// TODO: Change the assertion to equivalence when the above TODO is resolved.
	if !strings.Contains(out.String(), errorMessage) {
		t.Errorf(
			"Expected error message to be '%s', got '%s' instead",
			errorMessage, out.String())
	}
}

func Test_DirectoryList_getDetailsText_DoesNotReturnOutputFromPreviousCall(t *testing.T) {
	var seedDirectories []*mock.FileNode
	seedDirNamePart := "test"
	seedDirCount := 3

	seedDirectories = mock.GenerateSeedDirectories(seedDirNamePart, seedDirCount) // creates /test0, /test1, /test2
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 2, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	details := CreateDetailsView()

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), details, dirCtrl, nil)

	result0 := list.getDetailsText("test0")
	result1 := list.getDetailsText("test1")
	result2 := list.getDetailsText("test2")

	if strings.Contains(result1, result0) {
		t.Errorf(
			"Expected the first result not to contain the output of the second.\nFirst output:\n%s\nSecond output:\n%s\n",
			result0, result1)
	}

	if strings.Contains(result2, result1) {
		t.Errorf(
			"Expected the second result not to contain the output of the third.\nSecond output:\n%s\nThird output:\n%s\n",
			result1, result2)
	}
}

func getDirectoryControllerWithMockCommands(fileSystem *mock.FileSystem) *dirctrl.DefaultDirectoryController {
	dirCtrl := dirctrl.NewDefaultDirectoryController()
	dirCtrl.Commands = mock.NewDirectoryCommandsForVirtualFileSystem(fileSystem)

	return dirCtrl
}

func getAppWithDisabledExitHandlersAndOutputStreams(screen tcell.SimulationScreen) *App {
	app := NewApp(screen, io.Discard, io.Discard)
	app.handleNormalExit = func() {
		// Do nothing for test
	}
	app.handleErrorExit = func() {
		// Do nothing for test
	}

	return app
}

func Test_DirectoryList_handleDetailsInputCapture_SetsFocusToListWhenTabPressed(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	focus := CreateFilterForm()
	list := CreateDirectoryList(app, nil, focus, nil, nil, nil, nil)
	app.SetFocus(focus)

	list.handleDetailsInputCapture(tcell.NewEventKey(tcell.KeyTab, rune(tcell.KeyTab), tcell.ModNone))

	expected := list
	result := app.GetFocus()

	if result != expected {
		t.Errorf("Expected object '%v', got '%v' instead", expected, result)
	}
}

func Test_DirectoryList_handleDetailsInputCapture_AppExitsWhenShortcutKeyIsPressed(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	list := CreateDirectoryList(app, nil, nil, nil, nil, nil, nil)
	exited := false
	app.handleNormalExit = func() {
		exited = true
	}

	list.handleDetailsInputCapture(tcell.NewEventKey('q', 'q', tcell.ModNone))

	if !exited {
		t.Error("Expected app to exit, but it did not")
	}
}

func Test_DirectoryList_handleDetailsInputCapture_DoesNotReturnEventForHandledKeyPresses(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	list := CreateDirectoryList(app, nil, nil, nil, nil, nil, nil)
	app.handleNormalExit = func() {
		// Do nothing for test
	}

	handledKeyPresses := map[tcell.Key]rune{
		tcell.KeyEscape: rune(tcell.KeyEscape),
		tcell.KeyEnter:  rune(tcell.KeyEnter),
		tcell.KeyTab:    rune(tcell.KeyTab),
		'q':             'q',
	}

	for i := range handledKeyPresses {
		result := list.handleDetailsInputCapture(tcell.NewEventKey(i, handledKeyPresses[i], tcell.ModNone))

		if result != nil {
			t.Errorf("Did not expect the tcell.EventKey '%v' from the key '%c'",
				result, handledKeyPresses[i])
		}

	}
}

func Test_DirectoryList_handleFilterEntry_SetsListTitleWhenFilterEntered(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	filter := CreateFilterForm()
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	mockFileSystem := mock.NewMockFileSystem(nil, 2, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	list := CreateDirectoryList(app, tview.NewTextView(), filter, tview.NewPages(), CreateDetailsView(), dirCtrl, nil)

	filter.filterMethod.SetCurrentOption(filterMethodGlobPattern)

	filterText := "bananas"
	expectedListTitle := fmt.Sprintf("%v - Filter: %v", listTitle, filterText)

	filter.SetText(filterText)
	list.handleFilterEntry(tcell.KeyEnter)

	result := list.GetTitle()

	if result != expectedListTitle {
		t.Errorf("Expected list title to be '%s', got '%s' instead", expectedListTitle, result)
	}
}

func Test_DirectoryList_handleFilterEntry_ResetsListTitleToDefaultWhenFilterIsEmpty(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	filter := CreateFilterForm()
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	mockFileSystem := mock.NewMockFileSystem(nil, 2, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	list := CreateDirectoryList(app, tview.NewTextView(), filter, tview.NewPages(), CreateDetailsView(), dirCtrl, nil)

	filterText := ""
	expectedListTitle := listTitle

	filter.SetText(filterText)
	list.handleFilterEntry(tcell.KeyEnter)

	result := list.GetTitle()

	if result != expectedListTitle {
		t.Errorf("Expected list title to be '%s', got '%s' instead", expectedListTitle, result)
	}
}

func Test_DirectoryList_handleFilterEntry_SetsAppFocusToListWhenFilterEntered(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	filter := CreateFilterForm()
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	mockFileSystem := mock.NewMockFileSystem(nil, 2, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	list := CreateDirectoryList(app, tview.NewTextView(), filter, tview.NewPages(), CreateDetailsView(), dirCtrl, nil)
	app.SetFocus(filter)

	list.handleFilterEntry(tcell.KeyEnter)

	expectedFocus := list
	result := app.GetFocus()

	if result != expectedFocus {
		t.Errorf("Expected object '%v', got '%v' instead", expectedFocus, result)
	}
}

func Test_DirectoryList_handleFilterEntry_PerformsFilterWhenListReloaded(t *testing.T) {
	seedDirectories := mock.GetSampleExampleSeedDirectories()

	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 1, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	filter := CreateFilterForm()
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	list := CreateDirectoryList(app, tview.NewTextView(), filter, tview.NewPages(), CreateDetailsView(), dirCtrl, nil)

	filter.SetText("example2")

	expectedFail := fmt.Sprintf(
		"sample0%s sample1%s example0%s example1%s",
		dirctrl.OsPathSeparator,
		dirctrl.OsPathSeparator,
		dirctrl.OsPathSeparator,
		dirctrl.OsPathSeparator)

	list.handleFilterEntry(tcell.KeyEnter)

	foundFilteredItem := false
	for i := 0; i < list.GetItemCount(); i++ {
		itemText, _ := list.GetItemText(i)

		if strings.Contains(expectedFail, itemText) {
			t.Errorf("Expected a filtered list but it contained the invalid item '%s'", itemText)
		}

		if itemText == "example2" {
			foundFilteredItem = true
		}
	}

	if !foundFilteredItem {
		t.Errorf("Expected to find the filtered item 'example2', but it was not present")
	}
}

func Test_DirectoryList_handleFilterEntry_ClearsFilterTextWhenFilterEntered(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	filter := CreateFilterForm()
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	mockFileSystem := mock.NewMockFileSystem(nil, 2, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	list := CreateDirectoryList(app, tview.NewTextView(), filter, tview.NewPages(), CreateDetailsView(), dirCtrl, nil)

	filterText := "bananas"
	expectedFilterText := ""

	filter.SetText(filterText)
	list.handleFilterEntry(tcell.KeyEnter)

	result := filter.GetText()

	if result != expectedFilterText {
		t.Errorf("Expected filter text to be '%s', got '%s' instead",
			expectedFilterText, result)
	}
}

func Test_DirectoryList_handleFilterEntry_HidesFilterPaneAfterEntry(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	filter := CreateFilterForm()
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	pages := tview.NewPages()
	mockFileSystem := mock.NewMockFileSystem(nil, 2, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	list := CreateDirectoryList(app, tview.NewTextView(), filter, pages, CreateDetailsView(), dirCtrl, nil)

	pages.AddPage("List", list, false, true)
	pages.AddPage("Filter", filter, false, true)

	app.SetFocus(filter)

	frontPage, _ := pages.GetFrontPage()

	if frontPage != "Filter" {
		t.Error("This test requires that the 'Filter' page is the front page before the filter handler is called")
	}

	list.handleFilterEntry(tcell.KeyEnter)

	frontPage, _ = pages.GetFrontPage()

	if frontPage != "List" {
		t.Errorf("Expected 'List' to be the front page, got '%s' instead", frontPage)
	}
}

func Test_DirectoryList_handleFilterEntry_DoesNotApplyFilterIfEscIsPressed(t *testing.T) {
	seedDirectories := mock.GetSampleExampleSeedDirectories()

	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 1, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	filter := CreateFilterForm()
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	list := CreateDirectoryList(app, tview.NewTextView(), filter, tview.NewPages(), CreateDetailsView(), dirCtrl, nil)

	filter.SetText("example2")

	list.handleFilterEntry(tcell.KeyEscape)

	var result string
	for i := 0; i < list.GetItemCount(); i++ {
		itemText, _ := list.GetItemText(i)
		itemText = strings.TrimRight(itemText, dirctrl.OsPathSeparator)

		if !list.isMenuItem(itemText) {
			result = result + " " + itemText
		}
	}

	for _, dir := range seedDirectories {
		if !strings.Contains(result, dir.File.Name()) {
			t.Errorf("Expected an unfiltered list but it did not contain item '%s'", dir.File.Name())
		}
	}
}

func Test_DirectoryList_handleScrollPosition_SetsVerticalOffsetToMinimumWhenFirstItemSelected(t *testing.T) {
	var seedDirectories []*mock.FileNode
	seedDirNamePart := "test"
	seedDirCount := 20

	seedDirectories = mock.GenerateSeedDirectories(seedDirNamePart, seedDirCount)
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 2, 10)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)

	listWidth := 20
	listHeight := 10
	list.SetRect(0, 0, listWidth, listHeight)

	if _, err := mockFileSystem.Cd("/"); err != nil {
		t.Fatal(err)
	}
	list.load()

	expectedVPosition := 0
	list.SetCurrentItem(expectedVPosition)
	list.SetOffset(1, 0)
	vPosition, _ := list.handleScrollPosition()

	if vPosition != expectedVPosition {
		t.Errorf("Expected vertical scroll position to be %v, got '%v' instead",
			expectedVPosition, vPosition)
	}
}

func Test_DirectoryList_handleScrollPosition_SetsVerticalOffsetToMaximumWhenLastItemSelected(t *testing.T) {
	var seedDirectories []*mock.FileNode
	seedDirNamePart := "test"
	seedDirCount := 20

	seedDirectories = mock.GenerateSeedDirectories(seedDirNamePart, seedDirCount)
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 2, 10)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)

	listWidth := 20
	listHeight := 10
	list.SetRect(0, 0, listWidth, listHeight)
	_, _, _, listPageHeight := list.GetInnerRect()

	if _, err := mockFileSystem.Cd("/"); err != nil {
		t.Fatal(err)
	}
	list.load()

	itemCount := list.GetItemCount()
	expectedVPosition := itemCount - listPageHeight
	list.SetCurrentItem(itemCount - 1)
	list.SetOffset(listHeight, 0)
	vPosition, _ := list.handleScrollPosition()

	if vPosition != expectedVPosition {
		t.Errorf("Expected vertical scroll position to be %v, got '%v' instead",
			expectedVPosition, vPosition)
	}
}

func Test_DirectoryList_handleScrollPosition_SetsScrollPositionAsListOffsetWhenNeitherFirstNorLastItemSelected(t *testing.T) {
	var seedDirectories []*mock.FileNode
	seedDirNamePart := "test"
	seedDirCount := 20

	seedDirectories = mock.GenerateSeedDirectories(seedDirNamePart, seedDirCount)
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 2, 10)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)

	listWidth := 20
	listHeight := 10
	list.SetRect(0, 0, listWidth, listHeight)

	if _, err := mockFileSystem.Cd("/"); err != nil {
		t.Fatal(err)
	}
	list.load()

	expectedVPosition := 1
	list.SetCurrentItem(5)
	list.SetOffset(expectedVPosition, 0)
	vPosition, _ := list.handleScrollPosition()

	if vPosition != expectedVPosition {
		t.Errorf("Expected vertical scroll position to be %v, got '%v' instead",
			expectedVPosition, vPosition)
	}
}

func Test_DirectoryList_handleInputCapture_LeftArrowKeyNavigatesToPreviousDirectory(t *testing.T) {
	seedDirectories := mock.GetHierarchicalSeedDirectories()
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 4, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	if _, err := mockFileSystem.Cd("/testA/testB"); err != nil {
		t.Fatal(err)
	}

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)
	if err := initializeCurrentDirectoryForTest(list); err != nil {
		t.Fatal(err)
	}
	list.load()

	list.handleInputCapture(tcell.NewEventKey(tcell.KeyLeft, rune(tcell.KeyLeft), tcell.ModNone))

	expectedCurrentDir := mock.NormalizePath("/testA")

	if list.currentDir != expectedCurrentDir {
		t.Errorf("Expected the current directory to be '%s', got '%s' instead", expectedCurrentDir, list.currentDir)
	}
}

// initializeCurrentDirectoryForTest sets the current directory for the DirectoryList without having to call DirectoryList.Init().
// This ensures that we run tests with the minimum necessary setup.
func initializeCurrentDirectoryForTest(list *DirectoryList) error {
	var err error
	list.currentDir, err = list.dirUtil.GetInitialDirectory()

	return err
}

func Test_DirectoryList_handleInputCapture_RightArrowKeyNavigatesToNextDirectory(t *testing.T) {
	seedDirectories := mock.GetHierarchicalSeedDirectories()
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 4, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	if _, err := mockFileSystem.Cd("/testA/"); err != nil {
		t.Fatal(err)
	}

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)
	if err := initializeCurrentDirectoryForTest(list); err != nil {
		t.Fatal(err)
	}
	list.load()

	setSelectedItem(list, "testB")

	list.handleInputCapture(tcell.NewEventKey(tcell.KeyRight, rune(tcell.KeyRight), tcell.ModNone))

	expectedCurrentDir := mock.NormalizePath("/testA/testB")

	if list.currentDir != expectedCurrentDir {
		t.Errorf("Expected the current directory to be '%s', got '%s' instead", expectedCurrentDir, list.currentDir)
	}
}

func setSelectedItem(list *DirectoryList, itemText string) {
	isSelected := false
	for i := 0; i < list.GetItemCount() && !isSelected; i++ {
		item, _ := list.GetItemText(i)
		if item == itemText {
			list.SetCurrentItem(i)
			isSelected = true
		}
	}
}

func Test_DirectoryList_handleInputCapture_UpArrowKeyDisplaysDetailsForPreviousItem(t *testing.T) {
	var seedDirectories []*mock.FileNode
	seedDirNamePart := "test"
	seedDirCount := 3
	seedDirectories = mock.GenerateSeedDirectories(seedDirNamePart, seedDirCount)
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 2, 3)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	details := CreateDetailsView()
	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), details, dirCtrl, nil)

	expectedDetailsText := map[string]string{}

	for i := 0; i < seedDirCount; i++ {
		dirName := seedDirNamePart + strconv.Itoa(i)
		expectedDetailsText[dirName] = list.getDetailsText(dirName)
	}

	list.load()

	setSelectedItem(list, "test2")

	list.handleInputCapture(tcell.NewEventKey(tcell.KeyUp, rune(tcell.KeyUp), tcell.ModNone))

	result := details.GetText(false)

	if result != expectedDetailsText["test1"] {
		t.Errorf(
			"Expected the following details with '%s' selected:\n%s\nGot the following instead:\n%s\n",
			mock.NormalizePath("test1"), expectedDetailsText["test1"], result)
	}
}

func Test_DirectoryList_handleInputCapture_DownArrowKeyDisplaysDetailsForPreviousItem(t *testing.T) {
	var seedDirectories []*mock.FileNode
	seedDirNamePart := "test"
	seedDirCount := 3
	seedDirectories = mock.GenerateSeedDirectories(seedDirNamePart, seedDirCount)
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 2, 3)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	details := CreateDetailsView()
	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), details, dirCtrl, nil)

	expectedDetailsText := map[string]string{}

	for i := 0; i < seedDirCount; i++ {
		dirName := seedDirNamePart + strconv.Itoa(i)
		expectedDetailsText[dirName] = list.getDetailsText(dirName)
	}

	list.load()

	setSelectedItem(list, "test0")

	list.handleInputCapture(tcell.NewEventKey(tcell.KeyDown, rune(tcell.KeyDown), tcell.ModNone))

	result := details.GetText(false)

	if result != expectedDetailsText["test1"] {
		t.Errorf(
			"Expected the following details with '%s' selected:\n%s\nGot the following instead:\n%s\n",
			mock.NormalizePath("test1"), expectedDetailsText["test1"], result)
	}
}

func Test_DirectoryList_handleInputCapture_TabKeySetsFocusToDetailsPane(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	details := CreateDetailsView()
	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), details, nil, nil)

	app.SetFocus(list)

	list.handleInputCapture(tcell.NewEventKey(tcell.KeyTab, rune(tcell.KeyTab), tcell.ModNone))

	result := app.GetFocus()

	if result != details {
		t.Errorf("Expected the app focus to be '%v', got '%v' instead", details, result)
	}
}

func Test_DirectoryList_handleLeftKeyEvent_SetsCurrentDirectoryToOneLevelUpFromPreviousValue(t *testing.T) {
	seedDirectories := mock.GetHierarchicalSeedDirectories()
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 4, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	if _, err := mockFileSystem.Cd("/testA/testB"); err != nil {
		t.Fatal(err)
	}

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)
	if err := initializeCurrentDirectoryForTest(list); err != nil {
		t.Fatal(err)
	}
	list.load()

	list.handleLeftKeyEvent()

	expectedCurrentDir := mock.NormalizePath("/testA")

	if list.currentDir != expectedCurrentDir {
		t.Errorf("Expected the current directory to be '%s', got '%s' instead",
			expectedCurrentDir, list.currentDir)
	}
}

func Test_DirectoryList_handleLeftKeyEvent_SetsDirectoryListTitleWhenNavigating(t *testing.T) {
	seedDirectories := mock.GetHierarchicalSeedDirectories()
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 4, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	if _, err := mockFileSystem.Cd("/testA/testB"); err != nil {
		t.Fatal(err)
	}

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)
	if err := initializeCurrentDirectoryForTest(list); err != nil {
		t.Fatal(err)
	}
	list.load()
	list.SetTitle("This should not appear")

	list.handleLeftKeyEvent()

	expectedTitle := listTitle
	result := list.GetTitle()

	if result != expectedTitle {
		t.Errorf("Expected title to be '%s', got '%s' instead", expectedTitle, result)
	}
}

func Test_DirectoryList_handleLeftKeyEvent_DoesNotNavigateWhenAtRootDirectory(t *testing.T) {
	seedDirectories := mock.GetHierarchicalSeedDirectories()
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 4, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	details := CreateDetailsView()

	if _, err := mockFileSystem.Cd("/"); err != nil {
		t.Fatal(err)
	}

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), details, dirCtrl, nil)
	list.load()
	list.loadDetailsForCurrentDirectory()

	expectedDetails := list.getDetailsText(mock.NormalizePath("/"))
	expectedCurrentDir := list.currentDir

	list.handleLeftKeyEvent()

	resultDetails := details.GetText(false)
	resultCurrentDir := list.currentDir

	if resultDetails != expectedDetails {
		t.Errorf("Expected details to be the following:\n%s\nGot the following instead:\n%s\n",
			expectedDetails, resultDetails)
	}

	if resultCurrentDir != expectedCurrentDir {
		t.Errorf("Expected current directory to be '%s', got '%s' instead", expectedCurrentDir, resultCurrentDir)
	}
}

func Test_DirectoryList_handleLeftKeyEvent_ClearsFilterWhenNavigating(t *testing.T) {
	seedDirectories := mock.GetHierarchicalSeedDirectories()
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 4, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	if _, err := mockFileSystem.Cd("/testA/testB"); err != nil {
		t.Fatal(err)
	}

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)
	if err := initializeCurrentDirectoryForTest(list); err != nil {
		t.Fatal(err)
	}
	list.load()
	list.filterText = "This should not appear"

	list.handleLeftKeyEvent()

	expectedFilterText := ""
	result := list.filterText

	if result != expectedFilterText {
		t.Errorf("Expected title to be '%s', got '%s' instead", expectedFilterText, result)
	}
}

func Test_DirectoryList_handleLeftKeyEvent_LoadsListForParentDirectory(t *testing.T) {
	seedDirectories := mock.GetHierarchicalSeedDirectories()
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 4, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	if _, err := mockFileSystem.Cd("/testA"); err != nil {
		t.Fatal(err)
	}

	files, err := mockFileSystem.Ls(".")
	if err != nil {
		t.Fatal(err)
	}

	var expectedDirNames []string
	for _, file := range files {
		if file.IsDir() {
			expectedDirNames = append(expectedDirNames, file.Name())
		}
	}

	if _, err := mockFileSystem.Cd("/testA/testB"); err != nil {
		t.Fatal(err)
	}

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)
	if err := initializeCurrentDirectoryForTest(list); err != nil {
		t.Fatal(err)
	}
	list.load()

	list.handleLeftKeyEvent()

	for _, expected := range expectedDirNames {
		setSelectedItem(list, expected)

		selected, _ := list.GetItemText(list.GetCurrentItem())

		if expected != selected {
			t.Errorf("Expected '%s' to be in the list but it was not present", expected)
		}
	}
}

func Test_DirectoryList_handleLeftKeyEvent_LoadsListForRootDirectory(t *testing.T) {
	seedDirectories := mock.GetHierarchicalSeedDirectories()
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 4, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	files, err := mockFileSystem.Ls(".")
	if err != nil {
		t.Fatal(err)
	}

	var expectedDirNames []string
	for _, file := range files {
		if file.IsDir() {
			expectedDirNames = append(expectedDirNames, file.Name())
		}
	}

	if _, err := mockFileSystem.Cd("/testA"); err != nil {
		t.Fatal(err)
	}

	runTest := func(prependDriveLetter bool) {
		list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)
		if err := initializeCurrentDirectoryForTest(list); err != nil {
			t.Fatal(err)
		}
		list.load()

		if prependDriveLetter && runtime.GOOS == "windows" {
			list.currentDir = "C:" + list.currentDir
		}

		list.handleLeftKeyEvent()

		for _, expected := range expectedDirNames {
			setSelectedItem(list, expected)

			selected, _ := list.GetItemText(list.GetCurrentItem())

			if expected != selected {
				t.Errorf("Expected '%s' to be in the list but it was not present", expected)
			}
		}
	}

	runTest(false)
	runTest(true)
}

func Test_DirectoryList_handleLeftKeyEvent_ListsContentsOfDirectoryInDetailsPane(t *testing.T) {
	seedDirectories := mock.GetHierarchicalSeedDirectories()
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 4, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	details := CreateDetailsView()

	if _, err := mockFileSystem.Cd("/testA/testB"); err != nil {
		t.Fatal(err)
	}

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), details, dirCtrl, nil)
	if err := initializeCurrentDirectoryForTest(list); err != nil {
		t.Fatal(err)
	}
	list.load()

	expectedDetails := list.getDetailsText(mock.NormalizePath("/testA"))

	list.handleLeftKeyEvent()

	resultDetails := details.GetText(false)

	if resultDetails != expectedDetails {
		t.Errorf("Expected details to be the following:\n%s\nGot the following instead:\n%s\n",
			expectedDetails, resultDetails)
	}
}

func Test_DirectoryList_handleRightKeyEvent_SetsCurrentDirectoryToOneLevelDownFromPreviousValue(t *testing.T) {
	seedDirectories := mock.GetHierarchicalSeedDirectories()
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 4, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	if _, err := mockFileSystem.Cd("/testA/testB"); err != nil {
		t.Fatal(err)
	}

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)
	if err := initializeCurrentDirectoryForTest(list); err != nil {
		t.Fatal(err)
	}
	list.load()

	setSelectedItem(list, "testC")

	list.handleRightKeyEvent()

	expectedCurrentDir := mock.NormalizePath("/testA/testB/testC")

	if list.currentDir != expectedCurrentDir {
		t.Errorf("Expected the current directory to be '%s', got '%s' instead",
			expectedCurrentDir, list.currentDir)
	}
}

func Test_DirectoryList_handleRightKeyEvent_DoesNotNavigateIfListItemIsAMenuItem(t *testing.T) {
	mockFileSystem := mock.NewMockFileSystem(nil, 1, 10)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)
	list.load()

	setSelectedItem(list, listItemEnterDir)

	expectedCurrentDir := list.currentDir

	list.handleRightKeyEvent()

	resultCurrentDir := list.currentDir

	if resultCurrentDir != expectedCurrentDir {
		t.Errorf("Expected the current directory to be '%s', got '%s' instead",
			expectedCurrentDir, resultCurrentDir)
	}
}

func Test_DirectoryList_handleRightKeyEvent_ClearsFilterWhenNavigating(t *testing.T) {
	seedDirectories := mock.GetHierarchicalSeedDirectories()
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 4, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	if _, err := mockFileSystem.Cd("/testA"); err != nil {
		t.Fatal(err)
	}

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)
	list.load()
	list.filterText = "This should not appear"

	setSelectedItem(list, "testB")

	list.handleRightKeyEvent()

	expectedFilterText := ""
	result := list.filterText

	if result != expectedFilterText {
		t.Errorf("Expected filter text to be '%s', got '%s' instead", expectedFilterText, result)
	}
}

func Test_DirectoryList_handleRightKeyEvent_DoesNotNavigateIfDirectoryIsInaccessible(t *testing.T) {
	seedDirectories := mock.GetHierarchicalSeedDirectories()
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 4, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	if _, err := mockFileSystem.Cd("/testA"); err != nil {
		t.Fatal(err)
	}

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)
	list.load()

	setSelectedItem(list, "testB")

	expectedDetails := "[red]Directory inaccessible, unable to navigate. You may have insufficient privileges.[white]"
	expectedCurrentDir := list.currentDir

	dirCtrl.Commands = &mock.DirectoryCommands{
		ReadDirectoryFunc: func(dirname string) ([]fs.FileInfo, error) {
			return nil, errors.New("error triggered by test")
		},
		GetAbsolutePathFunc: func(path string) (string, error) {
			return "", errors.New("error triggered by test")
		},
		ScanDirectoryFunc: func(path string, callback func(dirName string)) error {
			return errors.New("error triggered by test")
		},
	}

	list.handleRightKeyEvent()

	resultDetails := list.details.GetText(false)
	resultCurrentDir := list.currentDir

	if resultDetails != expectedDetails {
		t.Errorf("Expected details to be the following:\n%s\nGot the following instead:\n%s\n",
			expectedDetails, resultDetails)
	}

	if resultCurrentDir != expectedCurrentDir {
		t.Errorf("Expected current directory to be '%s', got '%s' instead", expectedCurrentDir, resultCurrentDir)
	}
}

func Test_DirectoryList_handleRightKeyEvent_SetsDirectoryListTitleWhenNavigating(t *testing.T) {
	seedDirectories := mock.GetHierarchicalSeedDirectories()
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 4, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	if _, err := mockFileSystem.Cd("/testA/testB"); err != nil {
		t.Fatal(err)
	}

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)
	list.load()
	list.SetTitle("This should not appear")

	setSelectedItem(list, "testC")

	list.handleRightKeyEvent()

	expectedTitle := listTitle
	result := list.GetTitle()

	if result != expectedTitle {
		t.Errorf("Expected title to be '%s', got '%s' instead", expectedTitle, result)
	}
}

func Test_DirectoryList_handleRightKeyEvent_LoadsListForChildDirectory(t *testing.T) {
	seedDirectories := mock.GetHierarchicalSeedDirectories()
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 4, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	if _, err := mockFileSystem.Cd("/testA/testB"); err != nil {
		t.Fatal(err)
	}

	files, err := mockFileSystem.Ls(".")
	if err != nil {
		t.Fatal(err)
	}

	var expectedDirNames []string
	for _, file := range files {
		if file.IsDir() {
			expectedDirNames = append(expectedDirNames, file.Name())
		}
	}

	if _, err := mockFileSystem.Cd("/testA"); err != nil {
		t.Fatal(err)
	}

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)
	list.load()

	setSelectedItem(list, "testB")

	list.handleRightKeyEvent()

	for _, expected := range expectedDirNames {
		setSelectedItem(list, expected)

		selected, _ := list.GetItemText(list.GetCurrentItem())

		if expected != selected {
			t.Errorf("Expected '%s' to be in the list but it was not present", expected)
		}
	}
}

func Test_DirectoryList_handleRightKeyEvent_SetsCurrentDirectoryCorrectlyFromRootDirectory(t *testing.T) {
	seedDirectories := mock.GetHierarchicalSeedDirectories()
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 4, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	if _, err := mockFileSystem.Cd("/"); err != nil {
		t.Fatal(err)
	}

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)
	if err := initializeCurrentDirectoryForTest(list); err != nil {
		t.Fatal(err)
	}
	list.load()

	setSelectedItem(list, "testA")

	list.handleRightKeyEvent()

	expectedCurrentDir := mock.NormalizePath("/testA")

	if list.currentDir != expectedCurrentDir {
		t.Errorf("Expected the current directory to be '%s', got '%s' instead",
			expectedCurrentDir, list.currentDir)
	}
}

func Test_DirectoryList_load_LoadsChildDirectoriesOfCurrentDirectory(t *testing.T) {
	seedDirectories := mock.GetHierarchicalSeedDirectories()
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 4, 10)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	if _, err := mockFileSystem.Cd("/testA/testB"); err != nil {
		t.Fatal(err)
	}

	files, err := mockFileSystem.Ls(".")
	if err != nil {
		t.Fatal(err)
	}

	var expectedDirNames []string
	for _, file := range files {
		if file.IsDir() {
			expectedDirNames = append(expectedDirNames, file.Name())
		}
	}

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)

	list.load()

	for _, expectedDirName := range expectedDirNames {
		setSelectedItem(list, expectedDirName)
		result, _ := list.GetItemText(list.GetCurrentItem())

		if result != expectedDirName {
			t.Errorf("Expected the directory '%s' to be present in the list, but it was not", expectedDirName)
		}
	}
}

func Test_DirectoryList_load_SetsAppTitleBoxTextToCurrentDirectoryPath(t *testing.T) {
	seedDirectories := mock.GetHierarchicalSeedDirectories()
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 4, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	title := tview.NewTextView()

	expectedTitleText := mock.NormalizePath("/testA/testB/testC")
	if _, err := mockFileSystem.Cd("/testA/testB/testC"); err != nil {
		t.Fatal(err)
	}

	list := CreateDirectoryList(app, title, CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)
	if err := initializeCurrentDirectoryForTest(list); err != nil {
		t.Fatal(err)
	}

	list.load()

	result := title.GetText(true)

	if result != expectedTitleText {
		t.Errorf("Expected the text of the title box to be '%s', got '%s' instead", expectedTitleText, result)
	}
}

func Test_DirectoryList_load_LoadsChildDirectoriesOfRootDirectory(t *testing.T) {
	mockFileSystem := mock.NewMockFileSystem(nil, 1, 10)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)

	files, err := mockFileSystem.Ls(".")
	if err != nil {
		t.Fatal(err)
	}

	var expectedDirNames []string
	for _, file := range files {
		if file.IsDir() {
			expectedDirNames = append(expectedDirNames, file.Name())
		}
	}

	list.load()

	for _, expectedDirName := range expectedDirNames {
		setSelectedItem(list, expectedDirName)
		result, _ := list.GetItemText(list.GetCurrentItem())

		if result != expectedDirName {
			t.Errorf("Expected the directory '%s' to be present in the list, but it was not", expectedDirName)
		}
	}
}

func Test_DirectoryList_load_LoadsSymbolicLinkDirectories(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	tempDir, err := os.MkdirTemp("", uuid.NewString())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Fatal(err)
		}
	}()

	childTempDirs := map[string]string{
		"testA": "",
		"testB": "",
		"testC": "",
		"testD": "",
	}
	childTempDirs["testA"] = filepath.Join(tempDir, "testA")
	childTempDirs["testB"] = filepath.Join(childTempDirs["testA"], "testB")
	childTempDirs["testC"] = filepath.Join(childTempDirs["testB"], "testC")
	childTempDirs["testD"] = filepath.Join(childTempDirs["testC"], "testD")

	if err = os.MkdirAll(childTempDirs["testD"], fs.ModePerm); err != nil {
		t.Fatal(err)
	}

	createSymLink := func(linkname, canonPath string) {
		tempSymLink := filepath.Join(tempDir, linkname)
		if err = os.Symlink(canonPath, tempSymLink); err != nil {
			if runtime.GOOS == "windows" && strings.Contains(err.Error(), "A required privilege is not held by the client") {
				t.Skip("Test skipped due to insufficient privileges to run it")
			} else {
				t.Fatal(err)
			}
		}
	}

	createSymLink("testB", childTempDirs["testB"])
	createSymLink("testC", childTempDirs["testC"])
	createSymLink("testD", childTempDirs["testD"])

	// Ensure that symbolic links to files do not end up in the directory list
	tempFileName := "testfile.txt"
	tempFileNamePath := filepath.Join(childTempDirs["testB"], tempFileName)
	if testFile, err := os.Create(tempFileNamePath); err != nil {
		t.Fatal(err)
	} else if err = testFile.Close(); err != nil {
		t.Fatal(err)
	}
	createSymLink(tempFileName, tempFileNamePath)

	list := CreateDirectoryList(
		app,
		tview.NewTextView(),
		CreateFilterForm(),
		tview.NewPages(),
		CreateDetailsView(),
		dirctrl.NewDefaultDirectoryController(),
		nil)
	list.currentDir = tempDir
	list.load()

	for expectedDirName := range childTempDirs {
		setSelectedItem(list, expectedDirName)
		result, _ := list.GetItemText(list.GetCurrentItem())

		if result != expectedDirName {
			t.Errorf("Expected the directory '%s' to be present in the list, but it was not", expectedDirName)
		}
	}

	setSelectedItem(list, tempFileName)
	result, _ := list.GetItemText(list.GetCurrentItem())
	if result == tempFileName {
		t.Errorf("Expected the file '%s' not to be present in the list, but it was", tempFileName)
	}
}

func Test_DirectoryList_addNavigableItem_AddsListItemWhenFilterTextIsEmpty(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), nil, nil)

	expectedItem := "bananas"
	list.addNavigableItem(expectedItem)

	var allItems []string
	foundItem := false
	for i := 0; i < list.GetItemCount() && !foundItem; i++ {
		item, _ := list.GetItemText(i)
		foundItem = item == expectedItem
		allItems = append(allItems, item)
	}

	if !foundItem {
		t.Errorf("Expected to find '%s' in the list, got the following items instead:\n%s\n",
			expectedItem,
			fmt.Sprintf("%v", allItems))
	}
}

func Test_DirectoryList_addNavigableItem_AddsListItemWhenFilterTextMatchesDirectoryName(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), nil, nil)
	list.filterText = "bananas"

	expectedItem := "bananas"
	list.addNavigableItem(expectedItem)

	var allItems []string
	foundItem := false
	for i := 0; i < list.GetItemCount() && !foundItem; i++ {
		item, _ := list.GetItemText(i)
		foundItem = item == expectedItem
		allItems = append(allItems, item)
	}

	if !foundItem {
		t.Errorf("Expected to find '%s' in the list, got the following items instead:\n%s\n",
			expectedItem,
			fmt.Sprintf("%v", allItems))
	}
}

func Test_DirectoryList_addNavigableItem_AddsListItemWhenFilterTextMatchesGlobPattern(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), nil, nil)

	runTest := func(filter, expected string) {
		list.filterText = filter

		expectedItem := expected
		list.addNavigableItem(expectedItem)

		var allItems []string
		foundItem := false
		for i := 0; i < list.GetItemCount() && !foundItem; i++ {
			item, _ := list.GetItemText(i)
			foundItem = item == expectedItem
			allItems = append(allItems, item)
		}

		if !foundItem {
			t.Errorf("Expected to find '%s' in the list, got the following items instead:\n%s\n",
				expectedItem,
				fmt.Sprintf("%v", allItems))
		}
	}

	runTest("ban*", "bananas")
	runTest("*pples", "apples")
	runTest("*ang*", "oranges")
	runTest("*in*pl*s", "pineapples")
}

func Test_DirectoryList_getNavigableItemSelectionHandler_GeneratesFullPathToDirectory(t *testing.T) {
	var out bytes.Buffer
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	app.outputStream = &out
	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), nil, nil)

	currentDir := mock.NormalizePath("/oranges/apples")
	dirName := "bananas"
	list.currentDir = currentDir

	selectionHandler := list.getNavigableItemSelectionHandler(dirName)
	expectedPath := currentDir + dirctrl.OsPathSeparator + dirName

	selectionHandler()

	result := out.String()
	if result != expectedPath {
		t.Errorf("Expected the path to be '%s', got '%s' instead", expectedPath, result)
	}
}

func Test_DirectoryList_setDetailsText_ScrollsTextToTop(t *testing.T) {
	seedDirectories := mock.GetHierarchicalSeedDirectories()
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 4, 10)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	details := CreateDetailsView()

	if _, err := mockFileSystem.Cd("/testA/testB"); err != nil {
		t.Fatal(err)
	}

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), details, dirCtrl, nil)

	list.setDetailsText("testC")

	resultX, resultY := details.GetScrollOffset()
	if resultX != 0 || resultY != 0 {
		t.Errorf(
			"Expected the details pane scroll position to be reset, got the following coordinates instead: (%v, %v)",
			resultX, resultY)
	}
}

func Test_DirectoryList_setDetailsText_SetsTextToCurrentDirectoryItemsWhenDefaultListItemSelected(t *testing.T) {
	mockFileSystem := mock.NewMockFileSystem(nil, 2, 10)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	details := CreateDetailsView()

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), details, dirCtrl, nil)

	expectedDetailsText := list.getDetailsText(".")

	list.setDetailsText(listItemEnterDir)

	result := details.GetText(false)

	if result != expectedDetailsText {
		t.Errorf(
			"Expected details text to be the following:\n%s\nGot the following instead:\n%s\n",
			expectedDetailsText,
			result)
	}
}

func Test_DirectoryList_setDetailsText_SetsDetailsOfDirectoryListItem(t *testing.T) {
	var seedDirectories []*mock.FileNode
	seedDirNamePart := "test"
	seedDirCount := 5

	seedDirectories = mock.GenerateSeedDirectories(seedDirNamePart, seedDirCount)
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 2, 10)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	details := CreateDetailsView()

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), details, dirCtrl, nil)

	expectedDetailsText := map[string]string{}

	for i := 0; i < seedDirCount; i++ {
		dirName := seedDirNamePart + strconv.Itoa(i)
		expectedDetailsText[dirName] = list.getDetailsText(dirName)
	}

	app.SetFocus(list)
	list.load()

	testsExecuted := false

	for i := 0; i < list.List.GetItemCount(); i++ {
		itemText, _ := list.List.GetItemText(i)
		expected, isDir := expectedDetailsText[itemText]

		if isDir {
			testsExecuted = true
			list.setDetailsText(itemText)
			result := list.details.GetText(false)
			if result != expected {
				t.Errorf(
					"Expected directory '%s' to have the following details:\n%s\nGot the following instead:\n%s\n",
					itemText, expected, result)
			}
		}
	}

	if !testsExecuted {
		t.Error("No tests were run, test data may be invalid")
	}
}

func Test_DirectoryList_setDetailsText_SetsDetailsPaneTitleWhenHelpItemSelected(t *testing.T) {
	details := CreateDetailsView()

	list := CreateDirectoryList(nil, nil, nil, nil, details, nil, nil)

	list.setDetailsText(listItemHelp)

	resultTitle := details.GetTitle()
	if resultTitle != detailsHelpTitle {
		t.Errorf("Expected title to be '%s', got '%s' instead", detailsHelpTitle, resultTitle)
	}
}

func Test_DirectoryList_setDetailsText_SetsDetailsPaneTitleWhenHelpItemNotSelected(t *testing.T) {
	mockFileSystem := mock.NewMockFileSystem(nil, 2, 10)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	details := CreateDetailsView()

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), details, dirCtrl, nil).Init()

	list.setDetailsText(listItemHelp)
	list.setDetailsText(listItemEnterDir)

	resultTitle := details.GetTitle()
	if resultTitle != detailsViewTitle {
		t.Errorf("Expected title to be '%s', got '%s' instead", detailsViewTitle, resultTitle)
	}
}

func Test_DirectoryList_getNextItemIndex_GetsIndexOfNextListItemWhenNavigatingUp(t *testing.T) {
	list := CreateDirectoryList(nil, nil, nil, nil, nil, nil, nil)

	for i := 0; i < 5; i++ {
		list.AddItem("Test"+strconv.Itoa(i), "", 0, nil)
	}

	list.SetCurrentItem(1)

	expected := 0
	result := list.getNextItemIndex(false)

	if result != expected {
		t.Errorf("Expected the selected item index to be %v, got %v instead", expected, result)
	}

	list.SetCurrentItem(0)
	expected = 4
	result = list.getNextItemIndex(false)

	if result != expected {
		t.Errorf("Expected the selected item index to be %v, got %v instead", expected, result)
	}
}

func Test_DirectoryList_getNextItemIndex_GetsIndexOfNextListItemWhenNavigatingDown(t *testing.T) {
	list := CreateDirectoryList(nil, nil, nil, nil, nil, nil, nil)

	for i := 0; i < 5; i++ {
		list.AddItem("Test"+strconv.Itoa(i), "", 0, nil)
	}

	list.SetCurrentItem(1)

	expected := 2
	result := list.getNextItemIndex(true)

	if result != expected {
		t.Errorf("Expected the selected item index to be %v, got %v instead", expected, result)
	}

	list.SetCurrentItem(4)
	expected = 0
	result = list.getNextItemIndex(true)

	if result != expected {
		t.Errorf("Expected the selected item index to be %v, got %v instead", expected, result)
	}
}

func Test_DirectoryList_Init_SetsCurrentDirectoryCorrectly(t *testing.T) {
	seedDirectories := mock.GetHierarchicalSeedDirectories()
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 2, 10)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)

	expectedDir := mock.NormalizePath("/testA/testB")
	if _, err := mockFileSystem.Cd(expectedDir); err != nil {
		t.Fatal(err)
	}

	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), CreateDetailsView(), dirCtrl, nil)

	list.Init()

	if list.currentDir != expectedDir {
		t.Errorf("Expected the current directory to be '%s', got '%s' instead", expectedDir, list.currentDir)
	}
}

func Test_DirectoryList_handleHelpSelection_HelpShortcutSetsDetailsViewTitle(t *testing.T) {
	mockFileSystem := mock.NewMockFileSystem(nil, 2, 10)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	details := CreateDetailsView()

	list := CreateDirectoryList(app, tview.NewTextView(), CreateFilterForm(), tview.NewPages(), details, dirCtrl, nil).Init()

	list.handleHelpSelection()

	result := details.GetTitle()

	if result != detailsHelpTitle {
		t.Errorf("Expected details view title to be '%s', got '%s' instead", detailsHelpTitle, result)
	}
}