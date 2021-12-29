package utils

import (
	"bytes"
	"fmt"
	"github.com/karrick/godirwalk"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/tabwriter"
)

const OsPathSeparator = string(os.PathSeparator)

const (
	DirUnexpectedError = iota + 1
	DirUnprivilegedError
)

type DirectoryError struct {
	Err error
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
	buffer bytes.Buffer
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

type DefaultDirectoryCommands struct {}

func (*DefaultDirectoryCommands) ReadDirectory(dirname string) ([]fs.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

func (*DefaultDirectoryCommands) GetAbsolutePath(path string) (string, error) {
	return filepath.Abs(path)
}

func (*DefaultDirectoryCommands) ScanDirectory(
	path string,
	callback func(dirName string),
) error {
	// TODO: Determine if godirwalk can/should be used instead of ioutil.ReadDir in ReadDirectory,
	//  or vice-versa.
	scanner, err := godirwalk.NewScanner(path)
	if err != nil {
		return err
	}

	for scanner.Scan() {
		entry, err := scanner.Dirent()
		if err != nil {
			return err
		}

		if entry.IsDir() {
			callback(entry.Name())
		}
	}

	return nil
}

type DirectoryController interface {
	GetInitialDirectory() (string, error)
	DirectoryIsAccessible(dir string) bool
	GetDirectoryInfo(dir string) (string, error)
	DirectoryScanner
}

type DefaultDirectoryController struct {
	Writer InfoWriter
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
	// TODO: Ensure that the directory items are sorted in case-insensitive alphabetical order.
	files, err := d.Commands.ReadDirectory(dir)
	if err != nil {
		return "", &DirectoryError{
			Err: err,
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

func (d *DefaultDirectoryController) ScanDirectory(path string, callback func(dirName string)) error {
	return d.Commands.ScanDirectory(path, callback)
}
