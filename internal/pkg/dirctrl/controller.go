package dirctrl

import (
	"fmt"
)

// DirectoryController specifies the abstracted filesystem functions that ci uses.
type DirectoryController interface {
	GetInitialDirectory() (string, error)
	DirectoryIsAccessible(dir string) bool
	GetDirectoryInfo(dir string) (string, error)
	GetAbsolutePath(dir string) (string, error)
	ScanDirectory(path string, callback func(dirName string)) error
}

// DefaultDirectoryController contains a collection of methods that execute various
// commands on the filesystem.
type DefaultDirectoryController struct {
	Writer   InfoWriter
	Commands DirectoryCommands
}

// NewDefaultDirectoryController creates a new instance of DefaultDirectoryController with
// initialized methods.
func NewDefaultDirectoryController() *DefaultDirectoryController {
	return &DefaultDirectoryController{
		Writer:   NewDefaultInfoWriter(),
		Commands: &DefaultDirectoryCommands{},
	}
}

// GetInitialDirectory returns the path that ci was run from.
func (d *DefaultDirectoryController) GetInitialDirectory() (string, error) {
	return d.Commands.GetAbsolutePath(".")
}

// DirectoryIsAccessible determines if a directory is accessible to the current user.
func (d *DefaultDirectoryController) DirectoryIsAccessible(directory string) bool {
	_, err := d.Commands.ReadDirectory(directory)

	return err == nil
}

// GetDirectoryInfo returns a formatted list of files in the specified directory.
func (d *DefaultDirectoryController) GetDirectoryInfo(directory string) (string, error) {
	files, err := d.Commands.ReadDirectory(directory)
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

// GetAbsolutePath gets the full path of the specified directory.
func (d *DefaultDirectoryController) GetAbsolutePath(directory string) (string, error) {
	return d.Commands.GetAbsolutePath(directory)
}

// ScanDirectory iterates over each file in the path and executes a callback that is
// provided the name of that file.
func (d *DefaultDirectoryController) ScanDirectory(path string, callback func(dirName string)) error {
	return d.Commands.ScanDirectory(path, callback)
}
