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
	dir += utils.OsPathSeparator

	list, err := newDirectoryList(nil, nil, nil, nil, nil)
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

	list, err := newDirectoryList(nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	list.dirUtil = dirCtrl

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

	list, err := newDirectoryList(nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	list.dirUtil = dirCtrl

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

	list, err := newDirectoryList(app, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	list.dirUtil = dirCtrl

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
	seedFileNamePart := "test"
	seedFileCount := 3

	seedDirectories = generateSeedDirectories(seedFileNamePart, seedFileCount) // creates /test0, /test1, /test2
	mockFileSystem := tdUtils.NewMockFileSystem(seedDirectories, 2, 5)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitAndOutputStreams(screen)
	details := CreateDetailsPane()

	list, err := newDirectoryList(app, tview.NewTextView(), tview.NewInputField(), tview.NewPages(), details)
	if err != nil {
		t.Fatal(err)
	}
	list.dirUtil = dirCtrl

	if _, err := mockFileSystem.Cd("/test0"); err != nil {
		t.Fatal(err)
	}
	result0 := list.getDetailsText("test0")

	if _, err := mockFileSystem.Cd("/test1"); err != nil {
		t.Fatal(err)
	}
	result1 := list.getDetailsText("test1")

	if _, err := mockFileSystem.Cd("/test2"); err != nil {
		t.Fatal(err)
	}
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
			return fileSystem.Ls(), nil
		},
		getAbsolutePath: func(path string) (string, error) {
			return fileSystem.ReadLink(path)
		},
		scanDirectory: func(path string, callback func(dirName string)) error {
			for _, file := range fileSystem.Ls() {
				callback(file.Name())
			}
			return nil
		},
	}

	return dirCtrl
}

func getAppWithDisabledExitAndOutputStreams(screen tcell.SimulationScreen) *App {
	app := NewApp(screen, io.Discard, io.Discard)
	app.handleNormalExit = func() {
		// Do nothing for test
	}
	app.handleErrorExit = app.handleNormalExit

	return app
}


func Test_DirectoryList_getDetailsInputCaptureHandler_SetsFocusToList(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_getDetailsInputCaptureHandler_AppExitsWhenShortcutKeyIsPressed(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_getDetailsInputCaptureHandler_ReturnsEventIfNoKeyPressesHandled(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_getFilterEntryHandler_SetsListTitleWhenFilterEntered(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_getFilterEntryHandler_ResetsListTitleWhenFilterIsEmpty(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_getFilterEntryHandler_SetsAppFocusToListWhenFilterEntered(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_getFilterEntryHandler_PerformsFilterWhenListReloaded(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_getFilterEntryHandler_ClearsFilterTextWhenFilterEntered(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_getListScrollPositionHandler_SetsVerticalOffsetToMinimumWhenFirstItemSelected(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_getListScrollPositionHandler_SetsVerticalOffsetToMaximumWhenLastItemSelected(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_getInputCaptureHandler_LeftArrowKeyNavigatesToPreviousDirectory(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_getInputCaptureHandler_RightArrowKeyNavigatesToNextDirectory(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_getInputCaptureHandler_UpArrowKeyNavigatesToPreviousListItem(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_getInputCaptureHandler_DownArrowKeyNavigatesToNextListItem(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_getInputCaptureHandler_TabKeySetsFocusToDetailsPane(t *testing.T) {
	t.Error("Unimplemented test")
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

func Test_DirectoryList_handleRightKeyEvent_ListsContentsOfDirectoryInDetailsPane(t *testing.T) {
	t.Error("Unimplemented test")
}

func Test_DirectoryList_setDetailsText_ClearsOldTextBeforeDisplayingNewText(t *testing.T) {
	var seedDirectories []*tdUtils.MockFileNode
	seedFileNamePart := "test"
	seedFileCount := 5

	seedDirectories = generateSeedDirectories(seedFileNamePart, seedFileCount)
	mockFileSystem := tdUtils.NewMockFileSystem(seedDirectories, 2, 10)
	dirCtrl := getDirectoryControllerWithMockCommands(mockFileSystem)
	screen := tcell.NewSimulationScreen("")
	app := getAppWithDisabledExitAndOutputStreams(screen)
	details := CreateDetailsPane()

	list, err := newDirectoryList(app, tview.NewTextView(), tview.NewInputField(), tview.NewPages(), details)
	if err != nil {
		t.Fatal(err)
	}
	list.dirUtil = dirCtrl

	expectedDetailsText := map[string]string {}

	for i := 0; i < seedFileCount; i++ {
		dirName := seedFileNamePart + strconv.Itoa(i)
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
