package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/tabwriter"
)

var OsPathSeparator = string(os.PathSeparator)

func GetInitialDirectory() (string, error) {
	dir, err := filepath.Abs(".")

	return dir + OsPathSeparator, err
}

func DirectoryIsAccessible(dir string) bool {
	_, err := ioutil.ReadDir(dir)

	return err == nil
}

func GetDirectoryInfo(dir string) string {
	var out bytes.Buffer

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "[red]Unable to read directory details. You may have insufficient privileges.[white]"
	}

	writer := tabwriter.NewWriter(&out, 1, 2, 2, ' ', 0)

	// TODO: Create a function for printing each row of the tab output to reduce duplication.
	_, err = fmt.Fprintf(writer, "%v\t%v\t%v\t%v\n", "Mode", "Name", "ModTime", "Bytes")
	HandleError(err, true)

	_, err = fmt.Fprintf(writer, "%v\t%v\t%v\t%v\n", "----", "----", "-------", "-----")
	HandleError(err, true)

	for _, f := range files {
		dateFormat := "2006-01-02 3:04 PM"
		modTime := f.ModTime().Format(dateFormat)
		_, err := fmt.Fprintf(writer, "%v\t%v\t%v\t%v\n", f.Mode(), f.Name(), modTime, f.Size())
		HandleError(err, true)
	}

	HandleError(writer.Flush(), true)

	return out.String()
}