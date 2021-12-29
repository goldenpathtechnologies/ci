package ui

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/goldenpathtechnologies/ci/internal/pkg/utils"
	tdUtils "github.com/goldenpathtechnologies/ci/testdata/utils"
	"github.com/rivo/tview"
	"io"
	"io/fs"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

type MockDirectoryCommands struct {
	readDirectory func(dirname string) ([]fs.FileInfo, error)
	getAbsolutePath func(path string) (string, error)
	scanDirectory func(path string, callback func(dirName string)) error
}

func (m *MockDirectoryCommands) ReadDirectory(dirname string) ([]fs.FileInfo, error) {
	return m.readDirectory(dirname)
}

func (m *MockDirectoryCommands) GetAbsolutePath(path string) (string, error) {
	return m.getAbsolutePath(path)
}

func (m *MockDirectoryCommands) ScanDirectory(path string, callback func(dirName string)) error {
	return m.scanDirectory(path, callback)
}

type MockInfoWriter struct {
	write func(p []byte) (n int, err error)
	flush func() (string, error)
}

func (m *MockInfoWriter) Write(p []byte) (n int, err error) {
	return m.write(p)
}

func (m *MockInfoWriter) Flush() (string, error) {
	return m.flush()
}

func Test_DirectoryList_newDirectoryList_SetsCurrentDirectoryCorrectly(t *testing.T) {
	dir, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}

	list, err := newDirectoryList(nil, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	if list.currentDir != dir {
		t.Errorf("Expected the current directory to be '%s', got '%s' instead", dir, list.currentDir)
	}
}

func Test_DirectoryList_getDetailsText_ReturnsDirectoryDetails(t *testing.T) {
	var files []fs.FileInfo

	for i := 0; i < 10; i++ {
		files = append(files, tdUtils.GenerateMockFile())
	}

	dirCtrl := utils.NewDefaultDirectoryController()
	dirCtrl.Commands = &MockDirectoryCommands{
		readDirectory: func(dirname string) ([]fs.FileInfo, error) {
			return files, nil
		},
		getAbsolutePath: func(path string) (string, error) {
			return "", nil
		},
	}

	list, err := newDirectoryList(nil, nil, nil, nil, nil, dirCtrl)
	if err != nil {
		t.Fatal(err)
	}

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
	dirCtrl := utils.NewDefaultDirectoryController()
	dirCtrl.Commands = &MockDirectoryCommands{
		readDirectory: func(dirname string) ([]fs.FileInfo, error) {
			return nil, errors.New("unable to access directory")
		},
		getAbsolutePath: func(path string) (string, error) {
			return "", nil
		},
	}

	list, err := newDirectoryList(nil, nil, nil, nil, nil, dirCtrl)
	if err != nil {
		t.Fatal(err)
	}

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
		out bytes.Buffer
	)

	errorMessage := "error triggered by test"

	for i := 0; i < 10; i++ {
		files = append(files, tdUtils.GenerateMockFile())
	}

	dirCtrl := utils.NewDefaultDirectoryController()
	dirCtrl.Commands = &MockDirectoryCommands{
		readDirectory: func(dirname string) ([]fs.FileInfo, error) {
			return files, nil
		},
		getAbsolutePath: func(path string) (string, error) {
			return "", nil
		},
	}
	dirCtrl.Writer = &MockInfoWriter{
		write: func(p []byte) (n int, err error) {
			return 0, errors.New(errorMessage)
		},
		flush: func() (string, error) {
			return "", nil
		},
	}

	screen := tcell.NewSimulationScreen("") // "" = UTF-8 charset
	app := NewApp(screen, io.Discard, io.Discard)
	app.handleErrorExit = func() {
		// Do nothing for test
	}

	list, err := newDirectoryList(app, nil, nil, nil, nil, dirCtrl)
	if err != nil {
		t.Fatal(err)
	}

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
	var seedDirectories []*tdUtils.MockFileNode
	seedDirNamePart := "test"
	seedDirCount := 3

	seedDirectories = generateSeedDirectories(seedDirNamePart, seedDirCount) // creates /test0, /test1, /test2
	mockFileSystem := tdUtils.NewMockFileSystem(seedDirectories, 2, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	details := CreateDetailsPane()

	list, err := newDirectoryList(app, tview.NewTextView(), tview.NewInputField(), tview.NewPages(), details, dirCtrl)
	if err != nil {
		t.Fatal(err)
	}

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

func generateSeedDirectories(fileNamePrefix string, count int) []*tdUtils.MockFileNode {
	var directories []*tdUtils.MockFileNode

	for i := 0; i < count; i++ {
		directories = append(directories, &tdUtils.MockFileNode{
			File: tdUtils.MockFile{
				FileName:    fileNamePrefix + strconv.Itoa(i),
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: nil,
			Parent:   nil,
		})
	}

	return directories
}

func getDirectoryControllerWithMockCommands(fileSystem *tdUtils.MockFileSystem) *utils.DefaultDirectoryController {
	dirCtrl := utils.NewDefaultDirectoryController()
	dirCtrl.Commands = &MockDirectoryCommands{
		readDirectory: func(dirname string) ([]fs.FileInfo, error) {
			return fileSystem.Ls(dirname)
		},
		getAbsolutePath: func(path string) (string, error) {
			return fileSystem.ReadLink(path)
		},
		scanDirectory: func(path string, callback func(dirName string)) error {
			files, err := fileSystem.Ls(path)
			if err != nil {
				return err
			}
			for _, file := range files {
				callback(file.Name())
			}
			return nil
		},
	}

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


func Test_DirectoryList_getDetailsInputCaptureHandler_SetsFocusToListWhenTabPressed(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	focus := CreateFilterPane()
	list, err := newDirectoryList(app, nil, focus, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	app.SetFocus(focus)

	inputHandler := list.getDetailsInputCaptureHandler()

	inputHandler(tcell.NewEventKey(tcell.KeyTab, rune(tcell.KeyTab), tcell.ModNone))

	expected := list
	result := app.GetFocus()

	if result != expected {
		t.Errorf("Expected object '%v', got '%v' instead", expected, result)
	}
}

func Test_DirectoryList_getDetailsInputCaptureHandler_AppExitsWhenShortcutKeyIsPressed(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	list, err := newDirectoryList(app, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	exited := false
	app.handleNormalExit = func() {
		exited = true
	}

	inputHandler := list.getDetailsInputCaptureHandler()

	inputHandler(tcell.NewEventKey('q', 'q', tcell.ModNone))

	if !exited {
		t.Error("Expected app to exit, but it did not")
	}
}

func Test_DirectoryList_getDetailsInputCaptureHandler_DoesNotReturnEventForHandledKeyPresses(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	list, err := newDirectoryList(app, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	app.handleNormalExit = func() {
		// Do nothing for test
	}

	inputHandler := list.getDetailsInputCaptureHandler()

	handledKeyPresses := map[tcell.Key]rune{
		tcell.KeyEscape: rune(tcell.KeyEscape),
		tcell.KeyEnter: rune(tcell.KeyEnter),
		tcell.KeyTab: rune(tcell.KeyTab),
		'q': 'q',
	}

	for i := range handledKeyPresses {
		result := inputHandler(tcell.NewEventKey(i, handledKeyPresses[i], tcell.ModNone))

		if result != nil {
			t.Errorf("Did not expect the tcell.EventKey '%v' from the key '%c'",
				result, handledKeyPresses[i])
		}

	}
}

func Test_DirectoryList_getFilterEntryHandler_SetsListTitleWhenFilterEntered(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	filter := CreateFilterPane()
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	list, err := newDirectoryList(app, tview.NewTextView(), filter, tview.NewPages(), CreateDetailsPane(), nil)
	if err != nil {
		t.Fatal(err)
	}

	filterHandler := list.getFilterEntryHandler()
	filter.SetDoneFunc(filterHandler)

	filterText := "bananas"
	expectedListTitle := fmt.Sprintf("%v - Filter: %v", listUITitle, filterText)

	filter.SetText(filterText)
	filterHandler(tcell.KeyEnter)

	result := list.GetTitle()

	if result != expectedListTitle {
		t.Errorf("Expected list title to be '%s', got '%s' instead", expectedListTitle, result)
	}
}

func Test_DirectoryList_getFilterEntryHandler_ResetsListTitleToDefaultWhenFilterIsEmpty(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	filter := CreateFilterPane()
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	list, err := newDirectoryList(app, tview.NewTextView(), filter, tview.NewPages(), CreateDetailsPane(), nil)
	if err != nil {
		t.Fatal(err)
	}

	filterHandler := list.getFilterEntryHandler()
	filter.SetDoneFunc(filterHandler)

	filterText := ""
	expectedListTitle := listUITitle

	filter.SetText(filterText)
	filterHandler(tcell.KeyEnter)

	result := list.GetTitle()

	if result != expectedListTitle {
		t.Errorf("Expected list title to be '%s', got '%s' instead", expectedListTitle, result)
	}
}

func Test_DirectoryList_getFilterEntryHandler_SetsAppFocusToListWhenFilterEntered(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	filter := CreateFilterPane()
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	list, err := newDirectoryList(app, tview.NewTextView(), filter, tview.NewPages(), CreateDetailsPane(), nil)
	if err != nil {
		t.Fatal(err)
	}
	app.SetFocus(filter)

	filterHandler := list.getFilterEntryHandler()
	filter.SetDoneFunc(filterHandler)

	filterHandler(tcell.KeyEnter)

	expectedFocus := list
	result := app.GetFocus()

	if result != expectedFocus {
		t.Errorf("Expected object '%v', got '%v' instead", expectedFocus, result)
	}
}

func Test_DirectoryList_getFilterEntryHandler_PerformsFilterWhenListReloaded(t *testing.T) {
	seedDirectories := []*tdUtils.MockFileNode{
		{
			File:     tdUtils.MockFile{
				FileName:    "sample0",
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: nil,
			Parent:   nil,
		},
		{
			File:     tdUtils.MockFile{
				FileName:    "sample1",
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: nil,
			Parent:   nil,
		},
		{
			File:     tdUtils.MockFile{
				FileName:    "example0",
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: nil,
			Parent:   nil,
		},
		{
			File:     tdUtils.MockFile{
				FileName:    "example1",
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: nil,
			Parent:   nil,
		},
		{
			File:     tdUtils.MockFile{
				FileName:    "example2",
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: nil,
			Parent:   nil,
		},
	}

	mockFileSystem := tdUtils.NewMockFileSystem(seedDirectories, 1, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	filter := CreateFilterPane()
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	list, err := newDirectoryList(app, tview.NewTextView(), filter, tview.NewPages(), CreateDetailsPane(), dirCtrl)
	if err != nil {
		t.Fatal(err)
	}

	filterHandler := list.getFilterEntryHandler()
	filter.SetDoneFunc(filterHandler)
	filter.SetText("example2")

	expectedFail := fmt.Sprintf(
		"sample0%s sample1%s example0%s example1%s",
		utils.OsPathSeparator,
		utils.OsPathSeparator,
		utils.OsPathSeparator,
		utils.OsPathSeparator)

	filterHandler(tcell.KeyEnter)

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

func Test_DirectoryList_getFilterEntryHandler_ClearsFilterTextWhenFilterEntered(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	filter := CreateFilterPane()
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	list, err := newDirectoryList(app, tview.NewTextView(), filter, tview.NewPages(), CreateDetailsPane(), nil)
	if err != nil {
		t.Fatal(err)
	}

	filterHandler := list.getFilterEntryHandler()
	filter.SetDoneFunc(filterHandler)

	filterText := "bananas"
	expectedFilterText := ""

	filter.SetText(filterText)
	filterHandler(tcell.KeyEnter)

	result := filter.GetText()

	if result != expectedFilterText {
		t.Errorf("Expected filter text to be '%s', got '%s' instead",
			expectedFilterText, result)
	}
}

func Test_DirectoryList_getFilterEntryHandler_HidesFilterPaneAfterEntry(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	filter := CreateFilterPane()
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	pages := tview.NewPages()
	list, err := newDirectoryList(app, tview.NewTextView(), filter, pages, CreateDetailsPane(), nil)
	if err != nil {
		t.Fatal(err)
	}
	pages.AddPage("List", list, false, true)
	pages.AddPage("Filter", filter, false, true)

	app.SetFocus(filter)

	filterHandler := list.getFilterEntryHandler()
	filter.SetDoneFunc(filterHandler)

	frontPage, _ := pages.GetFrontPage()

	if frontPage != "Filter" {
		t.Error("This test requires that the 'Filter' page is the front page before the filter handler is called")
	}

	filterHandler(tcell.KeyEnter)

	frontPage, _ = pages.GetFrontPage()

	if frontPage != "List" {
		t.Errorf("Expected 'List' to be the front page, got '%s' instead", frontPage)
	}
}

func Test_DirectoryList_getFilterEntryHandler_DoesNotApplyFilterIfEscIsPressed(t *testing.T) {
	seedDirectories := []*tdUtils.MockFileNode{
		{
			File:     tdUtils.MockFile{
				FileName:    "sample0",
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: nil,
			Parent:   nil,
		},
		{
			File:     tdUtils.MockFile{
				FileName:    "sample1",
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: nil,
			Parent:   nil,
		},
		{
			File:     tdUtils.MockFile{
				FileName:    "example0",
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: nil,
			Parent:   nil,
		},
		{
			File:     tdUtils.MockFile{
				FileName:    "example1",
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: nil,
			Parent:   nil,
		},
		{
			File:     tdUtils.MockFile{
				FileName:    "example2",
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: nil,
			Parent:   nil,
		},
	}

	mockFileSystem := tdUtils.NewMockFileSystem(seedDirectories, 1, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	filter := CreateFilterPane()
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	list, err := newDirectoryList(app, tview.NewTextView(), filter, tview.NewPages(), CreateDetailsPane(), dirCtrl)
	if err != nil {
		t.Fatal(err)
	}

	filterHandler := list.getFilterEntryHandler()
	filter.SetDoneFunc(filterHandler)
	filter.SetText("example2")

	filterHandler(tcell.KeyEscape)

	var result string
	for i := 0; i < list.GetItemCount(); i++ {
		itemText, _ := list.GetItemText(i)
		itemText = strings.TrimRight(itemText, utils.OsPathSeparator)

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

func Test_DirectoryList_getScrollPositionHandler_SetsVerticalOffsetToMinimumWhenFirstItemSelected(t *testing.T) {
	var seedDirectories []*tdUtils.MockFileNode
	seedDirNamePart := "test"
	seedDirCount := 20

	seedDirectories = generateSeedDirectories(seedDirNamePart, seedDirCount)
	mockFileSystem := tdUtils.NewMockFileSystem(seedDirectories, 2, 10)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	list, err := newDirectoryList(app, tview.NewTextView(), tview.NewInputField(), tview.NewPages(), CreateDetailsPane(), dirCtrl)
	if err != nil {
		t.Fatal(err)
	}
	
	listWidth := 20
	listHeight := 10
	list.SetRect(0, 0, listWidth, listHeight)

	scrollPosHandler := list.getScrollPositionHandler()

	if _, err = mockFileSystem.Cd("/"); err != nil {
		t.Fatal(err)
	}
	list.load()

	expectedVPosition := 0
	list.SetCurrentItem(expectedVPosition)
	list.SetOffset(1, 0)
	vPosition, _ := scrollPosHandler()

	if vPosition != expectedVPosition {
		t.Errorf("Expected vertical scroll position to be %v, got '%v' instead",
			expectedVPosition, vPosition)
	}
}

func Test_DirectoryList_getScrollPositionHandler_SetsVerticalOffsetToMaximumWhenLastItemSelected(t *testing.T) {
	var seedDirectories []*tdUtils.MockFileNode
	seedDirNamePart := "test"
	seedDirCount := 20

	seedDirectories = generateSeedDirectories(seedDirNamePart, seedDirCount)
	mockFileSystem := tdUtils.NewMockFileSystem(seedDirectories, 2, 10)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	list, err := newDirectoryList(app, tview.NewTextView(), tview.NewInputField(), tview.NewPages(), CreateDetailsPane(), dirCtrl)
	if err != nil {
		t.Fatal(err)
	}
	
	listWidth := 20
	listHeight := 10
	list.SetRect(0, 0, listWidth, listHeight)
	_, _, _, listPageHeight := list.GetInnerRect()

	scrollPosHandler := list.getScrollPositionHandler()

	if _, err = mockFileSystem.Cd("/"); err != nil {
		t.Fatal(err)
	}
	list.load()

	itemCount := list.GetItemCount()
	expectedVPosition := itemCount - listPageHeight
	list.SetCurrentItem(itemCount-1)
	list.SetOffset(listHeight, 0)
	vPosition, _ := scrollPosHandler()

	if vPosition != expectedVPosition {
		t.Errorf("Expected vertical scroll position to be %v, got '%v' instead",
			expectedVPosition, vPosition)
	}
}

func Test_DirectoryList_getScrollPositionHandler_SetsScrollPositionAsListOffsetWhenNeitherFirstNorLastItemSelected(t *testing.T) {
	var seedDirectories []*tdUtils.MockFileNode
	seedDirNamePart := "test"
	seedDirCount := 20

	seedDirectories = generateSeedDirectories(seedDirNamePart, seedDirCount)
	mockFileSystem := tdUtils.NewMockFileSystem(seedDirectories, 2, 10)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	list, err := newDirectoryList(app, tview.NewTextView(), tview.NewInputField(), tview.NewPages(), CreateDetailsPane(), dirCtrl)
	if err != nil {
		t.Fatal(err)
	}
	
	listWidth := 20
	listHeight := 10
	list.SetRect(0, 0, listWidth, listHeight)

	scrollPosHandler := list.getScrollPositionHandler()

	if _, err = mockFileSystem.Cd("/"); err != nil {
		t.Fatal(err)
	}
	list.load()

	expectedVPosition := 1
	list.SetCurrentItem(5)
	list.SetOffset(expectedVPosition, 0)
	vPosition, _ := scrollPosHandler()

	if vPosition != expectedVPosition {
		t.Errorf("Expected vertical scroll position to be %v, got '%v' instead",
			expectedVPosition, vPosition)
	}
}

func Test_DirectoryList_getInputCaptureHandler_LeftArrowKeyNavigatesToPreviousDirectory(t *testing.T) {
	seedDirectories := []*tdUtils.MockFileNode{
		{
			File:     tdUtils.MockFile{
				FileName:    "testA",
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: []*tdUtils.MockFileNode{
				{
					File:     tdUtils.MockFile{
						FileName:    "testB",
						FileSize:    0,
						FileMode:    fs.ModeDir | fs.ModePerm,
						FileModTime: time.Now(),
					},
					Children: []*tdUtils.MockFileNode{
						{
							File:     tdUtils.MockFile{
								FileName:    "testC",
								FileSize:    0,
								FileMode:    fs.ModeDir | fs.ModePerm,
								FileModTime: time.Now(),
							},
							Children: nil,
							Parent:   nil,
						},
					},
					Parent:   nil,
				},
			},
			Parent:   nil,
		},
	}
	mockFileSystem := tdUtils.NewMockFileSystem(seedDirectories, 4, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	if _, err := mockFileSystem.Cd("/testA/testB"); err != nil {
		t.Fatal(err)
	}

	list, err := newDirectoryList(app, tview.NewTextView(), tview.NewInputField(), tview.NewPages(), CreateDetailsPane(), dirCtrl)
	if err != nil {
		t.Fatal(err)
	}
	list.load()

	inputHandler := list.getInputCaptureHandler()

	inputHandler(tcell.NewEventKey(tcell.KeyLeft, rune(tcell.KeyLeft), tcell.ModNone))

	expectedCurrentDir := tdUtils.NormalizePath("/testA")

	if list.currentDir != expectedCurrentDir {
		t.Errorf("Expected the current directory to be '%s', got '%s' instead", expectedCurrentDir, list.currentDir)
	}
}

func Test_DirectoryList_getInputCaptureHandler_RightArrowKeyNavigatesToNextDirectory(t *testing.T) {
	seedDirectories := []*tdUtils.MockFileNode{
		{
			File:     tdUtils.MockFile{
				FileName:    "testA",
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: []*tdUtils.MockFileNode{
				{
					File:     tdUtils.MockFile{
						FileName:    "testB",
						FileSize:    0,
						FileMode:    fs.ModeDir | fs.ModePerm,
						FileModTime: time.Now(),
					},
					Children: []*tdUtils.MockFileNode{
						{
							File:     tdUtils.MockFile{
								FileName:    "testC",
								FileSize:    0,
								FileMode:    fs.ModeDir | fs.ModePerm,
								FileModTime: time.Now(),
							},
							Children: nil,
							Parent:   nil,
						},
					},
					Parent:   nil,
				},
			},
			Parent:   nil,
		},
	}
	mockFileSystem := tdUtils.NewMockFileSystem(seedDirectories, 4, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	if _, err := mockFileSystem.Cd("/testA/"); err != nil {
		t.Fatal(err)
	}

	list, err := newDirectoryList(app, tview.NewTextView(), tview.NewInputField(), tview.NewPages(), CreateDetailsPane(), dirCtrl)
	if err != nil {
		t.Fatal(err)
	}
	list.load()

	inputHandler := list.getInputCaptureHandler()

	setSelectedItem(list, "testB")

	inputHandler(tcell.NewEventKey(tcell.KeyRight, rune(tcell.KeyRight), tcell.ModNone))

	expectedCurrentDir := tdUtils.NormalizePath("/testA/testB")

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

func Test_DirectoryList_getInputCaptureHandler_UpArrowKeyDisplaysDetailsForPreviousItem(t *testing.T) {
	var seedDirectories []*tdUtils.MockFileNode
	seedDirNamePart := "test"
	seedDirCount := 3
	seedDirectories = generateSeedDirectories(seedDirNamePart, seedDirCount)
	mockFileSystem := tdUtils.NewMockFileSystem(seedDirectories, 2, 3)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	details := CreateDetailsPane()
	list, err := newDirectoryList(app, tview.NewTextView(), tview.NewInputField(), tview.NewPages(), details, dirCtrl)
	if err != nil {
		t.Fatal(err)
	}

	expectedDetailsText := map[string]string {}

	for i := 0; i < seedDirCount; i++ {
		dirName := seedDirNamePart + strconv.Itoa(i)
		expectedDetailsText[dirName] = list.getDetailsText(dirName)
	}

	list.load()

	setSelectedItem(list, "test2")

	inputHandler := list.getInputCaptureHandler()
	inputHandler(tcell.NewEventKey(tcell.KeyUp, rune(tcell.KeyUp), tcell.ModNone))

	result := details.GetText(false)

	if result != expectedDetailsText["test1"] {
		t.Errorf(
			"Expected the following details with '%s' selected:\n%s\nGot the following instead:\n%s\n",
			tdUtils.NormalizePath("test1"), expectedDetailsText["test1"], result)
	}
}



func Test_DirectoryList_getInputCaptureHandler_DownArrowKeyDisplaysDetailsForPreviousItem(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_getInputCaptureHandler_TabKeySetsFocusToDetailsPane(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_handleLeftKeyEvent_SetsCurrentDirectoryToOneLevelUpFromPreviousValue(t *testing.T) {
	seedDirectories := []*tdUtils.MockFileNode{
		{
			File:     tdUtils.MockFile{
				FileName:    "testA",
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: []*tdUtils.MockFileNode{
				{
					File:     tdUtils.MockFile{
						FileName:    "testB",
						FileSize:    0,
						FileMode:    fs.ModeDir | fs.ModePerm,
						FileModTime: time.Now(),
					},
					Children: []*tdUtils.MockFileNode{
						{
							File:     tdUtils.MockFile{
								FileName:    "testC",
								FileSize:    0,
								FileMode:    fs.ModeDir | fs.ModePerm,
								FileModTime: time.Now(),
							},
							Children: nil,
							Parent:   nil,
						},
					},
					Parent:   nil,
				},
			},
			Parent:   nil,
		},
	}
	mockFileSystem := tdUtils.NewMockFileSystem(seedDirectories, 4, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	if _, err := mockFileSystem.Cd("/testA/testB"); err != nil {
		t.Fatal(err)
	}

	list, err := newDirectoryList(app, tview.NewTextView(), tview.NewInputField(), tview.NewPages(), CreateDetailsPane(), dirCtrl)
	if err != nil {
		t.Fatal(err)
	}
	list.load()

	list.handleLeftKeyEvent()

	expectedCurrentDir := tdUtils.NormalizePath("/testA")

	if list.currentDir != expectedCurrentDir {
		t.Errorf("Expected the current directory to be '%s', got '%s' instead",
			expectedCurrentDir, list.currentDir)
	}
}

func Test_DirectoryList_handleLeftKeyEvent_SetsTitleWhenNavigating(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_handleLeftKeyEvent_DoesNotNavigateWhenAtRootDirectory(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_handleLeftKeyEvent_ClearsFilterWhenNavigating(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_handleLeftKeyEvent_LoadsListForParentDirectory(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_handleLeftKeyEvent_ListsContentsOfDirectoryInDetailsPane(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_handleRightKeyEvent_SetsCurrentDirectoryToOneLevelDownFromPreviousValue(t *testing.T) {
	seedDirectories := []*tdUtils.MockFileNode{
		{
			File:     tdUtils.MockFile{
				FileName:    "testA",
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: []*tdUtils.MockFileNode{
				{
					File:     tdUtils.MockFile{
						FileName:    "testB",
						FileSize:    0,
						FileMode:    fs.ModeDir | fs.ModePerm,
						FileModTime: time.Now(),
					},
					Children: []*tdUtils.MockFileNode{
						{
							File:     tdUtils.MockFile{
								FileName:    "testC",
								FileSize:    0,
								FileMode:    fs.ModeDir | fs.ModePerm,
								FileModTime: time.Now(),
							},
							Children: nil,
							Parent:   nil,
						},
					},
					Parent:   nil,
				},
			},
			Parent:   nil,
		},
	}
	mockFileSystem := tdUtils.NewMockFileSystem(seedDirectories, 4, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)

	if _, err := mockFileSystem.Cd("/testA/testB"); err != nil {
		t.Fatal(err)
	}

	list, err := newDirectoryList(app, tview.NewTextView(), tview.NewInputField(), tview.NewPages(), CreateDetailsPane(), dirCtrl)
	if err != nil {
		t.Fatal(err)
	}
	list.load()

	setSelectedItem(list, tdUtils.NormalizePath("testC"))

	list.handleRightKeyEvent()

	expectedCurrentDir := tdUtils.NormalizePath("/testA/testB/testC")

	if list.currentDir != expectedCurrentDir {
		t.Errorf("Expected the current directory to be '%s', got '%s' instead",
			expectedCurrentDir, list.currentDir)
	}
}

func Test_DirectoryList_handleRightKeyEvent_DoesNotNavigateIfListItemIsAMenuItem(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_handleRightKeyEvent_ClearsFilterWhenNavigating(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_handleRightKeyEvent_DoesNotNavigateIfDirectoryIsInaccessible(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_handleRightKeyEvent_SetsTitleWhenNavigating(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_handleRightKeyEvent_LoadsListForChildDirectory(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_load_LoadsChildDirectoriesOfCurrentDirectory(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_load_SetsTitleToCurrentDirectoryPath(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_addNavigableItem_AddsItemToList(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	list, err := newDirectoryList(app, tview.NewTextView(), tview.NewInputField(), tview.NewPages(), CreateDetailsPane(), nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedItem := "bananas"
	list.addNavigableItem(expectedItem)

	var allItems []string
	foundItem := false
	for i := 0; i < list.GetItemCount() && !foundItem; i++ {
		item, _ := list.GetItemText(i)
		foundItem = item  == expectedItem
		allItems = append(allItems, item)
	}

	if !foundItem {
		t.Errorf("Expected to find '%s' in the list, got the following items instead:\n%s\n",
			expectedItem,
			fmt.Sprintf("%v", allItems))
	}
}

func Test_DirectoryList_addNavigableItem_AddsListItemWhenFilterTextIsEmpty(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_addNavigableItem_AddsListItemWhenFilterTextMatchesDirectoryName(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_addNavigableItem_AddsListItemWhenFilterTextMatchesGlobPattern(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_getNavigableItemSelectionHandler_GeneratesFullPathToDirectory(t *testing.T) {
	var out bytes.Buffer
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	app.outputStream = &out
	list, err := newDirectoryList(app, tview.NewTextView(), tview.NewInputField(), tview.NewPages(), CreateDetailsPane(), nil)
	if err != nil {
		t.Fatal(err)
	}
	currentDir := tdUtils.NormalizePath("/oranges/apples")
	dirName := "bananas"
	list.currentDir = currentDir

	selectionHandler := list.getNavigableItemSelectionHandler(dirName)
	expectedPath := currentDir + utils.OsPathSeparator + dirName

	selectionHandler()

	result := out.String()
	if result != expectedPath {
		t.Errorf("Expected the path to be '%s', got '%s' instead", expectedPath, result)
	}
}

func Test_DirectoryList_setDetailsText_ScrollsTextToTop(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_setDetailsText_SetsTextToCurrentDirectoryItemsWhenDefaultListItemSelected(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_setDetailsText_SetsTextToChildItemsOfSelectedListItem(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_setDetailsText_ClearsOldTextBeforeDisplayingNewText(t *testing.T) {
	var seedDirectories []*tdUtils.MockFileNode
	seedDirNamePart := "test"
	seedDirCount := 5

	seedDirectories = generateSeedDirectories(seedDirNamePart, seedDirCount)
	mockFileSystem := tdUtils.NewMockFileSystem(seedDirectories, 2, 10)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitHandlersAndOutputStreams(screen)
	details := CreateDetailsPane()

	list, err := newDirectoryList(app, tview.NewTextView(), tview.NewInputField(), tview.NewPages(), details, dirCtrl)
	if err != nil {
		t.Fatal(err)
	}

	expectedDetailsText := map[string]string {}

	for i := 0; i < seedDirCount; i++ {
		dirName := seedDirNamePart + strconv.Itoa(i)
		if _, err := mockFileSystem.Cd("/" + dirName); err != nil {
			t.Fatal(err)
		}
		expectedDetailsText[dirName] = list.getDetailsText(dirName)
	}

	if _, err := mockFileSystem.Cd("/"); err != nil {
		t.Fatal(err)
	}

	app.SetFocus(list)
	list.load()

	testsExecuted := false

	for i := 0; i < list.List.GetItemCount(); i++ {
		itemText, _ := list.List.GetItemText(i)
		itemText = strings.TrimRight(itemText, "/\\")
		expected, isDir := expectedDetailsText[itemText]

		if isDir {
			testsExecuted = true
			list.currentDir = itemText
			if _, err := mockFileSystem.Cd(itemText); err != nil {
				t.Fatal(err)
			}
			list.setDetailsText(true)
			result := list.details.GetText(false)
			if result != expected {
				t.Errorf(
					"Expected directory '%s' to have the following details:\n%s\nGot the following instead:\n%s\n",
					itemText, expected, result)
			}
		}

		if _, err := mockFileSystem.Cd("/"); err != nil {
			t.Fatal(err)
		}
	}

	if !testsExecuted {
		t.Error("No tests were run, test data may be invalid")
	}
}

func Test_DirectoryList_getNextItemIndex_GetsIndexOfNextListItemWhenNavigatingUp(t *testing.T) {
	var (
		list *DirectoryList
		err error
	)

	if list, err = newDirectoryList(nil, nil, nil, nil, nil, nil); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		list.AddItem("Test" + strconv.Itoa(i), "", 0, nil)
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
	var (
		list *DirectoryList
		err error
	)

	if list, err = newDirectoryList(nil, nil, nil, nil, nil, nil); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		list.AddItem("Test" + strconv.Itoa(i), "", 0, nil)
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
