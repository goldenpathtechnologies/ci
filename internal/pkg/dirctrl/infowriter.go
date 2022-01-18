package dirctrl

import (
	"bytes"
	"io"
	"text/tabwriter"
)

type InfoWriter interface {
	io.Writer
	Flush() (string, error)
}

type DefaultInfoWriter struct {
	buffer    bytes.Buffer
	tabWriter *tabwriter.Writer
}

func NewDefaultInfoWriter() *DefaultInfoWriter {
	infoWriter := &DefaultInfoWriter{}
	infoWriter.tabWriter = tabwriter.NewWriter(&infoWriter.buffer, 1, 2, 2, ' ', 0)

	return infoWriter
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
