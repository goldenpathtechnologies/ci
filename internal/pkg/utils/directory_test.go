package utils

import (
	"github.com/goldenpathtechnologies/ci/testdata/utils"
	"github.com/google/uuid"
	"io/fs"
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

	t.Logf("\nActual directory order: %v\n", actualDirOrder)
}

func Test_DefaultDirectoryCommands_getFileInfoSliceSortHandler_OrdersFileInfoItemWithPrefixedUnderscore(t *testing.T) {
	files := []fs.FileInfo{
		utils.MockFile{
			FileName:    "Somefile",
			FileSize:    0,
			FileMode:    0,
			FileModTime: time.Now(),
		},
		utils.MockFile{
			FileName:    "2ndfile",
			FileSize:    0,
			FileMode:    0,
			FileModTime: time.Now(),
		},
		utils.MockFile{
			FileName:    ".dotfile",
			FileSize:    0,
			FileMode:    0,
			FileModTime: time.Now(),
		},
		utils.MockFile{
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

	t.Logf("\nActual file order: %v\n", files)
}

func Test_DefaultDirectoryCommands_getFileInfoSliceSortHandler_SortIsCaseInsensitive(t *testing.T) {
	files := []fs.FileInfo{
		utils.MockFile{
			FileName:    "Someflies",
			FileSize:    0,
			FileMode:    0,
			FileModTime: time.Now(),
		},
		utils.MockFile{
			FileName:    "Zebras",
			FileSize:    0,
			FileMode:    0,
			FileModTime: time.Now(),
		},
		utils.MockFile{
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

	t.Logf("\nActual file order: %v\n", files)
}

func Test_DefaultDirectoryCommands_getFileInfoSliceSortHandler_SortsFileInfoItemWithEmptyName(t *testing.T) {
	files := []fs.FileInfo{
		utils.MockFile{
			FileName:    "",
			FileSize:    0,
			FileMode:    0,
			FileModTime: time.Now(),
		},
		utils.MockFile{
			FileName:    "file0",
			FileSize:    0,
			FileMode:    0,
			FileModTime: time.Now(),
		},
		utils.MockFile{
			FileName:    "file1",
			FileSize:    0,
			FileMode:    0,
			FileModTime: time.Now(),
		},
	}

	expectedFileOrder := []string{"file0", "file1", ""}

	commands := &DefaultDirectoryCommands{}

	sort.Slice(files, commands.getFileInfoSliceSortHandler(files))

	for i, expectedFileName := range expectedFileOrder {
		if files[i].Name() != expectedFileName {
			t.Errorf("Expected '%s' at position '%d', got '%s' instead", expectedFileName, i, files[i].Name())
		}
	}

	t.Logf("\nActual file order: %v\n", files)
}
