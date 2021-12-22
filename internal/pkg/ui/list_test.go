package ui

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/goldenpathtechnologies/ci/internal/pkg/utils"
	tdUtils "github.com/goldenpathtechnologies/ci/testdata/utils"
	"io"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"testing"
)

type MockDirectoryCommands struct {
	readDirectory func(dirname string) ([]fs.FileInfo, error)
	getAbsolutePath func(path string) (string, error)
}

func (m *MockDirectoryCommands) ReadDirectory(dirname string) ([]fs.FileInfo, error) {
	return m.readDirectory(dirname)
}

func (m *MockDirectoryCommands) GetAbsolutePath(path string) (string, error) {
	return m.getAbsolutePath(path)
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
		files = append(files, tdUtils.GenerateTestFile())
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
		files = append(files, tdUtils.GenerateTestFile())
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
		panic(err)
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


