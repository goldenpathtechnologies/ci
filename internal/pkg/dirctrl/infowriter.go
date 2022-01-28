package dirctrl

import (
	"bytes"
	"io"
	"text/tabwriter"
)

// InfoWriter is an extension of the default Writer interface. It is designed to
// abstract a wider variety of output methods.
type InfoWriter interface {
	io.Writer
	Flush() (string, error)
}

// DefaultInfoWriter is an InfoWriter wrapper of the tabwriter module that outputs
// tab formatted data to a buffer.
type DefaultInfoWriter struct {
	buffer    bytes.Buffer
	tabWriter *tabwriter.Writer
}

// NewDefaultInfoWriter creates a new instance of DefaultInfoWriter.
func NewDefaultInfoWriter() *DefaultInfoWriter {
	infoWriter := &DefaultInfoWriter{}
	infoWriter.tabWriter = tabwriter.NewWriter(&infoWriter.buffer, 1, 2, 2, ' ', 0)

	return infoWriter
}

// Write is an implementation of the io.Writer Write function that sends data to the
// underlying tabwriter output stream.
func (d *DefaultInfoWriter) Write(p []byte) (n int, err error) {
	return d.tabWriter.Write(p)
}

// Flush clears the DefaultInfoWriter buffers and returns its contents.
func (d *DefaultInfoWriter) Flush() (string, error) {
	err := d.tabWriter.Flush()
	output := d.buffer.String()
	d.buffer.Reset()
	return output, err
}
