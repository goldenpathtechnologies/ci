package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
)

const OsPathSeparator = string(os.PathSeparator)

const (
	DirUnexpectedError = iota + 1
	DirUnprivilegedError
)

type DirectoryError struct {
	Err       error
	ErrorCode int
}

func (d *DirectoryError) Error() string {
	return d.Err.Error()
}

type InfoWriter interface {
	io.Writer
	Flush() (string, error)
}

type DefaultInfoWriter struct {
	buffer    bytes.Buffer
	tabWriter *tabwriter.Writer
}

func (d *DefaultInfoWriter) Write(p []byte) (n int, err error) {
	return d.tabWriter.Write(p)
}

func (d *DefaultInfoWriter) Flush() (string, error) {
	err := d.tabWriter.Flush()
	output := d.buffer.String()
	d.buffer.Reset()
	return output, err
}

type DirectoryScanner interface {
	ScanDirectory(path string, callback func(dirName string)) error
}

type DirectoryCommands interface {
	ReadDirectory(dirname string) ([]fs.FileInfo, error)
	GetAbsolutePath(path string) (string, error)
	DirectoryScanner
}

type DefaultDirectoryCommands struct{}

func (d *DefaultDirectoryCommands) ReadDirectory(dirname string) ([]fs.FileInfo, error) {
	items, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, fmt.Errorf(
			"DefaultDirectoryCommands.ReadDirectory, unable to read directory: %w",
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
	files, err := d.ReadDirectory(path)

	if err != nil {
		return err
	}

	for _, file := range files {
		isSymLink := file.Mode() & fs.ModeSymlink != 0
		if file.IsDir() {
			callback(file.Name())
		} else if isSymLink {
			link, err := os.Readlink(filepath.Join(path, file.Name()))
			if err != nil {
				return nil
			}

			_, err = ioutil.ReadDir(link)
			if err == nil {
				callback(file.Name())
			}
		}
	}

	return nil
}

type DirectoryController interface {
	GetInitialDirectory() (string, error)
	DirectoryIsAccessible(dir string) bool
	GetDirectoryInfo(dir string) (string, error)
	GetAbsolutePath(dir string) (string, error)
	DirectoryScanner
}

type DefaultDirectoryController struct {
	Writer   InfoWriter
	Commands DirectoryCommands
}

func NewDefaultDirectoryController() *DefaultDirectoryController {
	infoWriter := &DefaultInfoWriter{}
	infoWriter.tabWriter = tabwriter.NewWriter(&infoWriter.buffer, 1, 2, 2, ' ', 0)

	return &DefaultDirectoryController{
		Writer:   infoWriter,
		Commands: &DefaultDirectoryCommands{},
	}
}

func (d *DefaultDirectoryController) GetInitialDirectory() (string, error) {
	return d.Commands.GetAbsolutePath(".")
}

func (d *DefaultDirectoryController) DirectoryIsAccessible(dir string) bool {
	_, err := d.Commands.ReadDirectory(dir)

	return err == nil
}

func (d *DefaultDirectoryController) GetDirectoryInfo(dir string) (string, error) {
	files, err := d.Commands.ReadDirectory(dir)
	if err != nil {
		return "", &DirectoryError{
			Err:       err,
			ErrorCode: DirUnprivilegedError,
		}
	}

	// TODO: Return with a message if there are no files in the directory.

	_, err = fmt.Fprintf(d.Writer, "%v\t%v\t%v\t%v\n", "Mode", "Name", "ModTime", "Bytes")
	if err != nil {
		return "", &DirectoryError{
			Err:       err,
			ErrorCode: DirUnexpectedError,
		}
	}

	_, err = fmt.Fprintf(d.Writer, "%v\t%v\t%v\t%v\n", "----", "----", "-------", "-----")
	if err != nil {
		return "", &DirectoryError{
			Err:       err,
			ErrorCode: DirUnexpectedError,
		}
	}

	for _, f := range files {
		dateFormat := "2006-01-02 3:04 PM"
		modTime := f.ModTime().Format(dateFormat)
		_, err = fmt.Fprintf(d.Writer, "%v\t%v\t%v\t%v\n", f.Mode(), f.Name(), modTime, f.Size())
		if err != nil {
			return "", &DirectoryError{
				Err:       err,
				ErrorCode: DirUnexpectedError,
			}
		}
	}

	var output string
	if output, err = d.Writer.Flush(); err != nil {
		return "", &DirectoryError{
			Err:       err,
			ErrorCode: DirUnexpectedError,
		}
	}

	// TODO: Prune last newline from output.
	return output, nil
}

func (d *DefaultDirectoryController) GetAbsolutePath(dir string) (string, error) {
	return d.Commands.GetAbsolutePath(dir)
}

func (d *DefaultDirectoryController) ScanDirectory(path string, callback func(dirName string)) error {
	return d.Commands.ScanDirectory(path, callback)
}
