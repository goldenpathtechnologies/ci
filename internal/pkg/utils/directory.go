package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/tabwriter"
)

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

var OsPathSeparator = string(os.PathSeparator)

func GetInitialDirectory() (string, error) {
	dir, err := filepath.Abs(".")

	return dir + OsPathSeparator, err
}

func DirectoryIsAccessible(dir string) bool {
	_, err := ioutil.ReadDir(dir)

	return err == nil
}

func GetDirectoryInfo(dir string) (string, error) {
	var out bytes.Buffer

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", &DirectoryError{
			Err: err,
			ErrorCode: DirUnprivilegedError,
		}
	}

	writer := tabwriter.NewWriter(&out, 1, 2, 2, ' ', 0)

	_, err = fmt.Fprintf(writer, "%v\t%v\t%v\t%v\n", "Mode", "Name", "ModTime", "Bytes")
	if err != nil {
		return "", &DirectoryError{
			Err:       err,
			ErrorCode: DirUnexpectedError,
		}
	}

	_, err = fmt.Fprintf(writer, "%v\t%v\t%v\t%v\n", "----", "----", "-------", "-----")
	if err != nil {
		return "", &DirectoryError{
			Err:       err,
			ErrorCode: DirUnexpectedError,
		}
	}

	for _, f := range files {
		dateFormat := "2006-01-02 3:04 PM"
		modTime := f.ModTime().Format(dateFormat)
		_, err = fmt.Fprintf(writer, "%v\t%v\t%v\t%v\n", f.Mode(), f.Name(), modTime, f.Size())
		if err != nil {
			return "", &DirectoryError{
				Err:       err,
				ErrorCode: DirUnexpectedError,
			}
		}
	}

	if err = writer.Flush(); err != nil {
		return "", &DirectoryError{
			Err:       err,
			ErrorCode: DirUnexpectedError,
		}
	}

	// TODO: Prune last newline from output.
	return out.String(), nil
}