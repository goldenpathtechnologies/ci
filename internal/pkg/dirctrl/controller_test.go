package dirctrl

import (
	"errors"
	"github.com/goldenpathtechnologies/ci/testdata/mock"
	"io/fs"
	"strings"
	"testing"
)

func Test_DefaultDirectoryController_GetDirectoryInfo_ReturnsErrorWhenDirectoryCanNotBeRead(t *testing.T) {
	expectedErrorMessage := "error triggered by test"
	dirCtrl := NewDefaultDirectoryController()
	dirCtrl.Commands = &mock.DirectoryCommands{
		ReadDirectoryFunc: func(dirname string) ([]fs.FileInfo, error) {
			return nil, errors.New(expectedErrorMessage)
		},
	}

	if _, err := dirCtrl.GetDirectoryInfo("."); err == nil {
		t.Error("Expected returned error not to be nil")
	} else if !strings.Contains(err.Error(), expectedErrorMessage) {
		t.Errorf("Expected error to contain '%s', got the following error instead:\n%s\n", expectedErrorMessage, err.Error())
	}
}

func Test_DefaultDirectoryController_GetDirectoryInfo_ReturnsErrorWhenUnableToOutputToInfoWriter(t *testing.T) {
	dirCtrl := NewDefaultDirectoryController()
	oldWriter := dirCtrl.Writer

	// numCallsBeforeError is a way to help target specific calls to the InfoWriter's Write() func. I'm
	// still unsure how I feel about this approach, but it's the best I can come up with so far to ensure
	// all statements in GetDirectoryInfo() are tested.
	numCallsBeforeError := 0

	expectedErrorMessage := "error triggered by test"
	dirCtrl.Writer = &mock.InfoWriter{
		WriteFunc: func(p []byte) (n int, err error) {
			if numCallsBeforeError <= 0 {
				return 0, errors.New(expectedErrorMessage)
			} else {
				numCallsBeforeError--
				return oldWriter.Write(p)
			}
		},
		FlushFunc: oldWriter.Flush,
	}

	const maxDirectories = 3
	seedDirectories := mock.GenerateSeedDirectories("test", maxDirectories)
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 1, maxDirectories)
	dirCtrl.Commands = mock.NewDirectoryCommandsForVirtualFileSystem(mockFileSystem)

	testErrorOutput := func(origErr error, dErr *DirectoryError, isDirUnexpectedError bool) {
		if !isDirUnexpectedError {
			t.Errorf(
				"Expected the error returned to be of type DirectoryError, got the following instead:\n%v\n",
				origErr)
		} else if dErr.ErrorCode != DirUnexpectedError {
			t.Errorf(
				"Expected an error code of '%d', got '%d' instead",
				DirUnexpectedError, dErr.ErrorCode)
		} else if !strings.Contains(dErr.Error(), expectedErrorMessage) {
			t.Errorf(
				"Expected the original error containing the message '%s' to be forwarded, got the following error message instead:\n%s\n",
				expectedErrorMessage, dErr.Error())
		}
	}

	// Test error during printing of column headers
	_, err := dirCtrl.GetDirectoryInfo(mock.NormalizePath("/"))
	dErr, isDirUnexpectedError := err.(*DirectoryError)
	testErrorOutput(err, dErr, isDirUnexpectedError)

	// Test error during printing of column header divider
	numCallsBeforeError = 1
	_, err = dirCtrl.GetDirectoryInfo(mock.NormalizePath("/"))
	dErr, isDirUnexpectedError = err.(*DirectoryError)
	testErrorOutput(err, dErr, isDirUnexpectedError)

	// Test error during printing of file info
	numCallsBeforeError = 2
	_, err = dirCtrl.GetDirectoryInfo(mock.NormalizePath("/"))
	dErr, isDirUnexpectedError = err.(*DirectoryError)
	testErrorOutput(err, dErr, isDirUnexpectedError)
}

func Test_DefaultDirectoryController_GetDirectoryInfo_ReturnsErrorWhenUnableToGenerateOutputFromInfoWriter(t *testing.T) {
	dirCtrl := NewDefaultDirectoryController()
	oldWriter := dirCtrl.Writer

	expectedErrorMessage := "error triggered by test"
	dirCtrl.Writer = &mock.InfoWriter{
		WriteFunc: oldWriter.Write,
		FlushFunc: func() (string, error) {
			return "", errors.New(expectedErrorMessage)
		},
	}

	mockFileSystem := mock.NewMockFileSystem(nil, 1, 5)
	dirCtrl.Commands = mock.NewDirectoryCommandsForVirtualFileSystem(mockFileSystem)

	_, err := dirCtrl.GetDirectoryInfo(mock.NormalizePath("/"))
	dErr, isDirUnexpectedError := err.(*DirectoryError)
	if !isDirUnexpectedError {
		t.Errorf(
			"Expected the error returned to be of type DirectoryError, got the following instead:\n%v\n",
			err)
	} else if dErr.ErrorCode != DirUnexpectedError {
		t.Errorf(
			"Expected an error code of '%d', got '%d' instead",
			DirUnexpectedError, dErr.ErrorCode)
	} else if !strings.Contains(dErr.Error(), expectedErrorMessage) {
		t.Errorf(
			"Expected the original error containing the message '%s' to be forwarded, got the following error message instead:\n%s\n",
			expectedErrorMessage, dErr.Error())
	}
}

func Test_DefaultDirectoryController_GetDirectoryInfo_ReturnsListOfFilesInDirectory(t *testing.T) {
	const maxDirectories = 3
	dirCtrl := NewDefaultDirectoryController()
	seedDirectories := mock.GenerateSeedDirectories("test", maxDirectories)
	mockFileSystem := mock.NewMockFileSystem(seedDirectories, 1, maxDirectories)
	dirCtrl.Commands = mock.NewDirectoryCommandsForVirtualFileSystem(mockFileSystem)

	output, err := dirCtrl.GetDirectoryInfo(mock.NormalizePath("/"))

	if err != nil {
		t.Fatal(err)
	}

	containsAllDirectories :=
		strings.Contains(output, "test0") &&
			strings.Contains(output, "test1") &&
			strings.Contains(output, "test2")
	if !containsAllDirectories {
		t.Errorf("Expected to see specified files in output, got the following output instead:\n%s\n", output)
	}
}
