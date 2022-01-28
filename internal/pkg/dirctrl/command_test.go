package dirctrl

import (
	"bytes"
	"fmt"
	"github.com/goldenpathtechnologies/ci/testdata/mock"
	"github.com/google/uuid"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"
)

func Test_DefaultDirectoryCommands_ScanDirectory_ScansDirectoryItemsInCaseInsensitiveAlphabeticalOrder(t *testing.T) {
	tempDir, err := os.MkdirTemp("", uuid.NewString())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Fatal(err)
		}
	}()

	nestedDirs := map[string]string{
		"testA": "",
		"testB": "",
		"testC": "",
		"testD": "",
	}
	nestedDirs["testA"] = filepath.Join(tempDir, "testA")
	nestedDirs["testB"] = filepath.Join(nestedDirs["testA"], "testB")
	nestedDirs["testC"] = filepath.Join(nestedDirs["testB"], "testC")
	nestedDirs["testD"] = filepath.Join(nestedDirs["testC"], "testD")

	if err = os.MkdirAll(nestedDirs["testD"], fs.ModePerm); err != nil {
		t.Fatal(err)
	}

	createSymLink := func(linkname string, canonPath string) {
		tempSymLink := filepath.Join(tempDir, linkname)
		if err = os.Symlink(canonPath, tempSymLink); err != nil {
			if runtime.GOOS == "windows" && strings.Contains(err.Error(), "A required privilege is not held by the client") {
				t.Skip("Test skipped due to insufficient privileges to run it")
			} else {
				t.Fatal(err)
			}
		}
	}

	for _, v := range []string{"testB", "testC", "testD"} {
		createSymLink(v, nestedDirs[v])
	}

	mkDir := func(name string) {
		if err = os.Mkdir(filepath.Join(tempDir, name), fs.ModePerm); err != nil {
			t.Fatal(err)
		}
	}

	directoryNames := []string{
		".test", "test", "Testing", "apples", "bananas", "zebras", "Zeppelin", "Gnome", "dwarf",
		"_test", "Numbers", "3Amigos", "7samurai", "0Regrets", "night", ".github",
	}

	for _, v := range directoryNames {
		mkDir(v)
	}

	expectedDirOrder := []string {
		".github", ".test", "_test", "0Regrets", "3Amigos", "7samurai", "apples", "bananas",
		"dwarf", "Gnome", "night", "Numbers", "test", "testA", "testB", "testC", "testD",
		"Testing", "zebras", "Zeppelin",
	}

	commands := &DefaultDirectoryCommands{}
	i := 0
	var actualDirOrder []string
	if err = commands.ScanDirectory(tempDir, func(dirName string) {
		if dirName != expectedDirOrder[i] {
			t.Errorf("Expected '%s' at position '%d', got '%s' instead", expectedDirOrder[i], i, dirName)
		}
		actualDirOrder = append(actualDirOrder, dirName)
		i++
	}); err != nil {
		t.Fatal(err)
	}

	if t.Failed() {
		t.Logf("\nActual directory order: %v\n", actualDirOrder)
	}
}

func Test_DefaultDirectoryCommands_ScanDirectory_ReturnsErrorWhenProvidedWithNilCallback(t *testing.T) {
	tempDir, err := os.MkdirTemp("", uuid.NewString())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Fatal(err)
		}
	}()
	createDir := getCreateDirectoryForTestHandler(tempDir, t)
	createDir("test")

	commands := &DefaultDirectoryCommands{}
	err = commands.ScanDirectory(tempDir, nil)

	if err == nil {
		t.Fatal("Expected to return an error when provided with a nil callback")
	}

	expectedErrorMessage := "callback function must not be nil"

	if !strings.Contains(err.Error(), expectedErrorMessage) {
		t.Errorf("Expected to see the error message '%s' in the following error output:\n%s\n", expectedErrorMessage, err.Error())
	}
}

func getCreateDirectoryForTestHandler(tempDir string, t *testing.T) func(dirname string) {
	// TODO: This function and other similar ones in this file could potentially be moved to file.go
	//  in the testdata/mock directory. Do so once this function is needed in an additional suite.
	return func(dirname string) {
		if err := os.MkdirAll(filepath.Join(tempDir, dirname), fs.ModePerm); err != nil {
			t.Fatal(err)
		}
	}
}

func Test_DefaultDirectoryCommands_ScanDirectory_ReturnsAnnotatedErrorOnInvalidDirectoryName(t *testing.T) {
	tempDir, err := os.MkdirTemp("", uuid.NewString())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Fatal(err)
		}
	}()

	commands := &DefaultDirectoryCommands{}
	err = commands.ScanDirectory(filepath.Join(tempDir, "nonexistent"), func(dirName string) {})

	if err == nil {
		t.Fatal("Expected to return an error when provided with an invalid path")
	}

	expectedErrorMessage := "unable to scan directory, the path is invalid"

	if !strings.Contains(err.Error(), expectedErrorMessage) {
		t.Errorf("Expected to see the error message '%s' in the following error output:\n%s\n", expectedErrorMessage, err.Error())
	}
}

func Test_DefaultDirectoryCommands_ScanDirectory_LogsOccurrencesOfInvalidSymlinksWithoutRunningCallbackFunction(t *testing.T) {
	tempDir, err := os.MkdirTemp("", uuid.NewString())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Fatal(err)
		}
	}()

	createDir := getCreateDirectoryForTestHandler(tempDir, t)
	createDir("test")

	createSymlink := getCreateSymlinkForTestHandler(tempDir, t)
	createSymlink("testlink", filepath.Join(tempDir, "test"))

	if err = os.Remove(filepath.Join(tempDir, "test")); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	log.SetOutput(&out)

	commands := &DefaultDirectoryCommands{}
	callbackExecuted := false
	err = commands.ScanDirectory(tempDir, func(dirName string) {
		callbackExecuted = true
	})

	expectedLogMessage := fmt.Sprintf("the symlink '%s' is invalid", filepath.Join(tempDir, "testlink"))

	if !strings.Contains(out.String(), expectedLogMessage) {
		t.Errorf("Expected to see the message '%s' in the following log output:\n%s\n", expectedLogMessage, out.String())
	}

	if callbackExecuted {
		t.Error("Expected the callback function not to be executed")
	}
}

func getCreateFileForTestHandler(tempDir string, t *testing.T) func(filename string) {
	// TODO: This function and other similar ones in this file could potentially be moved to file.go
	//  in the testdata/mock directory. Do so once this function is needed in an additional suite.
	return func(filename string) {
		if testFile, err := os.Create(filepath.Join(tempDir, filename)); err != nil {
			t.Fatal(err)
		} else if err = testFile.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

func getCreateSymlinkForTestHandler(tempDir string, t *testing.T) func(linkname, canonPath string) {
	// TODO: This function and other similar ones in this file could potentially be moved to file.go
	//  in the testdata/mock directory. Do so once this function is needed in an additional suite.
	return func(linkname string, canonPath string) {
		tempSymLink := filepath.Join(tempDir, linkname)
		if err := os.Symlink(canonPath, tempSymLink); err != nil {
			if runtime.GOOS == "windows" && strings.Contains(err.Error(), "A required privilege is not held by the client") {
				t.Skip("Test skipped due to insufficient privileges to run it")
			} else {
				t.Fatal(err)
			}
		}
	}
}

func Test_DefaultDirectoryCommands_getFileInfoSliceSortHandler_OrdersFileInfoItemWithPrefixedUnderscore(t *testing.T) {
	files := []fs.FileInfo{
		mock.File{
			FileName:    "Somefile",
			FileSize:    0,
			FileMode:    0,
			FileModTime: time.Now(),
		},
		mock.File{
			FileName:    "2ndfile",
			FileSize:    0,
			FileMode:    0,
			FileModTime: time.Now(),
		},
		mock.File{
			FileName:    ".dotfile",
			FileSize:    0,
			FileMode:    0,
			FileModTime: time.Now(),
		},
		mock.File{
			FileName:    "_underscorefile",
			FileSize:    0,
			FileMode:    0,
			FileModTime: time.Now(),
		},
	}

	expectedFileOrder := []string{".dotfile", "_underscorefile", "2ndfile", "Somefile"}

	commands := &DefaultDirectoryCommands{}

	sort.Slice(files, commands.getFileInfoSliceSortHandler(files))

	for i, expectedFileName := range expectedFileOrder {
		if files[i].Name() != expectedFileName {
			t.Errorf("Expected '%s' at position '%d', got '%s' instead", expectedFileName, i, files[i].Name())
		}
	}

	if t.Failed() {
		t.Logf("\nActual file order: %v\n", files)
	}
}

func Test_DefaultDirectoryCommands_getFileInfoSliceSortHandler_SortIsCaseInsensitive(t *testing.T) {
	files := []fs.FileInfo{
		mock.File{
			FileName:    "Someflies",
			FileSize:    0,
			FileMode:    0,
			FileModTime: time.Now(),
		},
		mock.File{
			FileName:    "Zebras",
			FileSize:    0,
			FileMode:    0,
			FileModTime: time.Now(),
		},
		mock.File{
			FileName:    "somefile",
			FileSize:    0,
			FileMode:    0,
			FileModTime: time.Now(),
		},
	}

	expectedFileOrder := []string{"somefile", "Someflies", "Zebras"}

	commands := &DefaultDirectoryCommands{}

	sort.Slice(files, commands.getFileInfoSliceSortHandler(files))

	for i, expectedFileName := range expectedFileOrder {
		if files[i].Name() != expectedFileName {
			t.Errorf("Expected '%s' at position '%d', got '%s' instead", expectedFileName, i, files[i].Name())
		}
	}

	if t.Failed() {
		t.Logf("\nActual file order: %v\n", files)
	}
}

func Test_DefaultDirectoryCommands_getFileInfoSliceSortHandler_SortsFileInfoItemWithEmptyName(t *testing.T) {
	files := []fs.FileInfo{
		mock.File{
			FileName:    "",
			FileSize:    0,
			FileMode:    0,
			FileModTime: time.Now(),
		},
		mock.File{
			FileName:    "file0",
			FileSize:    0,
			FileMode:    0,
			FileModTime: time.Now(),
		},
		mock.File{
			FileName:    "file1",
			FileSize:    0,
			FileMode:    0,
			FileModTime: time.Now(),
		},
		mock.File{
			FileName:    "",
			FileSize:    0,
			FileMode:    0,
			FileModTime: time.Now(),
		},
	}

	expectedFileOrder := []string{"file0", "file1", "", ""}

	commands := &DefaultDirectoryCommands{}

	sort.Slice(files, commands.getFileInfoSliceSortHandler(files))

	for i, expectedFileName := range expectedFileOrder {
		if files[i].Name() != expectedFileName {
			t.Errorf("Expected '%s' at position '%d', got '%s' instead", expectedFileName, i, files[i].Name())
		}
	}

	if t.Failed() {
		t.Logf("\nActual file order: %v\n", files)
	}
}

func Test_DefaultDirectoryCommands_ReadDirectory_ReturnsErrorOnInvalidDirectoryName(t *testing.T) {
	tempDir, err := os.MkdirTemp("", uuid.NewString())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Fatal(err)
		}
	}()

	commands := &DefaultDirectoryCommands{}

	invalidPath := filepath.Join(tempDir, "nonexistent")
	_, err = commands.ReadDirectory(invalidPath)

	if err == nil {
		t.Fatal("Expected to return a non-nil error when provided an invalid path")
	}

	expectedErrorMessage := "unable to read directory"

	if !strings.Contains(err.Error(), expectedErrorMessage) {
		t.Errorf("Expected to see the error message '%s' in the following error output:\n%s\n", expectedErrorMessage, err.Error())
	}
}

func Test_DefaultDirectoryCommands_ReadDirectory_ReturnsSortedListOfFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", uuid.NewString())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Fatal(err)
		}
	}()

	createFile := getCreateFileForTestHandler(tempDir, t)

	createFile("zygote")
	createFile("manpower")
	createFile("rodent")
	createFile("helicopter")
	createFile("corrosion")

	expectedFileOrder := []string{"corrosion", "helicopter", "manpower", "rodent", "zygote"}

	commands := &DefaultDirectoryCommands{}

	if files, err := commands.ReadDirectory(tempDir); err != nil {
		t.Fatal(err)
	} else {
		for i, expectedFileName := range expectedFileOrder {
			if files[i].Name() != expectedFileName {
				t.Errorf("Expected '%s' at position '%d', got '%s' instead", expectedFileName, i, files[i].Name())
			}
		}
	}
}
