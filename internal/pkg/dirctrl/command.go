package dirctrl

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DirectoryCommands specifies the filesystem path functions that ci uses.
type DirectoryCommands interface {
	ReadDirectory(dirname string) ([]fs.FileInfo, error)
	GetAbsolutePath(path string) (string, error)
	ScanDirectory(path string, callback func(dirName string)) error
}

// DefaultDirectoryCommands is a placeholder struct for implemented methods
// of the DirectoryCommands interface.
type DefaultDirectoryCommands struct{}

// ReadDirectory returns a list of fs.FileInfo objects from the specified directory.
func (d *DefaultDirectoryCommands) ReadDirectory(dirname string) ([]fs.FileInfo, error) {
	// TODO: Consider switching to os.ReadDir(). Note that this func returns []fs.DirEntry.
	items, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to read directory: %w",
			err)
	}
	sort.Slice(items, d.getFileInfoSliceSortHandler(items))

	return items, nil
}

// getFileInfoSliceSortHandler sorts a list of files case-insensitively and with underscore
// prefixed files ordered above alphanumeric ones.
func (*DefaultDirectoryCommands) getFileInfoSliceSortHandler(items []fs.FileInfo) func (i, j int) bool {
	return func(i, j int) bool {
		compareI := items[i].Name()
		compareJ := items[j].Name()

		// Sort files with empty names to the end of the list (unlikely outside test environments)
		if len(compareI) == 0 {
			return false
		}
		if len(compareJ) == 0 {
			return true
		}

		runeI := rune(compareI[0])
		runeJ := rune(compareJ[0])

		// Sort files beginning with an underscore after dotfiles but before everything else
		if runeI == '_' && runeJ == '.' {
			return false
		} else if runeI == '_' && runeJ != '_' {
			return true
		} else if runeJ == '_' && runeI == '.' {
			return true
		} else if runeJ == '_' && runeI != '_' {
			return false
		}

		return strings.ToLower(compareI) < strings.ToLower(compareJ)
	}
}

// GetAbsolutePath gets the full path of the specified directory.
func (*DefaultDirectoryCommands) GetAbsolutePath(path string) (string, error) {
	return filepath.Abs(path)
}

// ScanDirectory iterates over each file in the path and executes a callback that is
// provided the name of that file.
func (d *DefaultDirectoryCommands) ScanDirectory(
	path string,
	callback func(dirName string),
) error {
	if callback == nil {
		return errors.New("callback function must not be nil")
	}

	files, err := d.ReadDirectory(path)

	if err != nil {
		return fmt.Errorf("unable to scan directory, the path is invalid: %w", err)
	}

	for _, file := range files {
		isSymLink := file.Mode() & fs.ModeSymlink != 0
		if file.IsDir() {
			callback(file.Name())
		} else if isSymLink {
			linkPath := filepath.Join(path, file.Name())
			link, err := os.Readlink(linkPath)
			if err != nil {
				log.Printf(
					fmt.Errorf(
						"an unexpected problem occured while reading symlink '%s': %w",
						linkPath,
						err).
						Error())
			}

			_, err = ioutil.ReadDir(link)
			if err != nil {
				log.Printf("the symlink '%s' is invalid", linkPath)
			} else {
				callback(file.Name())
			}
		}
	}

	return nil
}
