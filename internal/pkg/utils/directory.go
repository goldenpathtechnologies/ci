package utils

import (
	"bytes"
	"fmt"
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
	return d.buffer.String(), err
}

type DirectoryCommands interface {
	ReadDirectory(dirname string) ([]fs.FileInfo, error)
	GetAbsolutePath(path string) (string, error)
}

type DefaultDirectoryCommands struct {}

func (*DefaultDirectoryCommands) ReadDirectory(dirname string) ([]fs.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

func (*DefaultDirectoryCommands) GetAbsolutePath(path string) (string, error) {
	return filepath.Abs(path)
}

type DirectoryController interface {
	GetInitialDirectory() (string, error)
	DirectoryIsAccessible(dir string) bool
	GetDirectoryInfo(dir string) (string, error)
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
	dir, err := d.Commands.GetAbsolutePath(".")

	return dir + OsPathSeparator, err
}

func (d *DefaultDirectoryController) DirectoryIsAccessible(dir string) bool {
	_, err := d.Commands.ReadDirectory(dir)

	return err == nil
}

func (d *DefaultDirectoryController) GetDirectoryInfo(dir string) (string, error) {
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