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

type DirectoryCommands interface {
	ReadDirectory(dirname string) ([]fs.FileInfo, error)
	GetAbsolutePath(path string) (string, error)
	ScanDirectory(path string, callback func(dirName string)) error
}

type DefaultDirectoryCommands struct{}

func (d *DefaultDirectoryCommands) ReadDirectory(dirname string) ([]fs.FileInfo, error) {
	items, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to read directory: %w",
			err)
	}
	sort.Slice(items, d.getFileInfoSliceSortHandler(items))

	return items, nil
}

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

func (*DefaultDirectoryCommands) GetAbsolutePath(path string) (string, error) {
	return filepath.Abs(path)
}

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
