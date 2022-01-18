package dirctrl

import (
	"fmt"
	"os"
)

const OsPathSeparator = string(os.PathSeparator)

type DirectoryController interface {
	GetInitialDirectory() (string, error)
	DirectoryIsAccessible(dir string) bool
	GetDirectoryInfo(dir string) (string, error)
	GetAbsolutePath(dir string) (string, error)
	ScanDirectory(path string, callback func(dirName string)) error
}

type DefaultDirectoryController struct {
	Writer   InfoWriter
	Commands DirectoryCommands
}

func NewDefaultDirectoryController() *DefaultDirectoryController {
	return &DefaultDirectoryController{
		Writer:   NewDefaultInfoWriter(),
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
