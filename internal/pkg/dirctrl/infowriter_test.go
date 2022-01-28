package dirctrl

import (
	"fmt"
	"testing"
)

func Test_DefaultInfoWriter_Flush_ResetsAndReturnsContentsOfInternalBuffer(t *testing.T) {
	writer := NewDefaultInfoWriter()

	_, err := fmt.Fprintf(writer, "%s\t%s\t%s\n", "Test", "Test", "Test")
	if err != nil {
		t.Fatal(err)
	}

	_, err = writer.Flush()
	if err != nil {
		t.Fatal(err)
	}

	bufferContents := writer.buffer.String()
	if len(bufferContents) > 0 {
		t.Fatal("Expected the info writer buffer to be empty")
	}
}
