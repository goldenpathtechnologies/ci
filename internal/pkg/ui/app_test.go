package ui

import (
	"bytes"
	"errors"
	"github.com/gdamore/tcell/v2"
	"log"
	"testing"
)

const (
	testBufferEntrySequence = "\033[?1049h"
	testBufferExitSequence = "\033[?1049l"
)

func Test_App_enterScreenBuffer(t *testing.T) {
	var out bytes.Buffer
	screen := tcell.NewSimulationScreen("") // "" = UTF-8 charset
	app := NewApp(screen, &out, nil)
	app.screenBufferActive = false

	app.enterScreenBuffer()
	expected := testBufferEntrySequence
	result := out.String()

	if result != expected {
		t.Errorf("Expected hex output to be '%x', got '%x' instead", expected, result)
	}

	if !app.screenBufferActive {
		t.Errorf("Expected the screen buffer flag to be set")
	}
}

func Test_App_enterScreenBuffer_NoOutputWhenBufferActive(t *testing.T) {
	var out bytes.Buffer
	screen := tcell.NewSimulationScreen("")
	app := NewApp(screen, &out, nil)
	app.screenBufferActive = true

	app.enterScreenBuffer()
	expected := ""
	result := out.String()

	if result != expected {
		t.Errorf("Expected hex output to be '%x', got '%x' instead", expected, result)
	}
}

func Test_App_exitScreenBuffer(t *testing.T) {
	var out bytes.Buffer
	screen := tcell.NewSimulationScreen("")
	app := NewApp(screen, &out, nil)
	app.screenBufferActive = true

	app.exitScreenBuffer()
	expected := testBufferExitSequence
	result := out.String()

	if result != expected {
		t.Errorf("Expected hex output to be '%x', got '%x' instead", expected, result)
	}

	if app.screenBufferActive {
		t.Errorf("Expected the screen buffer flag to be unset")
	}
}

func Test_App_exitScreenBuffer_NoOutputWhenBufferInactive(t *testing.T) {
	var out bytes.Buffer
	screen := tcell.NewSimulationScreen("")
	app := NewApp(screen, &out, nil)
	app.screenBufferActive = false

	app.exitScreenBuffer()
	expected := ""
	result := out.String()

	if result != expected {
		t.Errorf("Expected hex output to be '%x', got '%x' instead", expected, result)
	}
}

func Test_App_PrintAndExit_PrintsToOutputStreamAndExits(t *testing.T) {
	var out bytes.Buffer
	screen := tcell.NewSimulationScreen("")
	app := NewApp(screen, &out, nil)
	performedExit := false
	app.handleNormalExit = func() {
		performedExit = true
	}

	expected := "bananas"
	app.PrintAndExit(expected)
	result := out.String()

	if result != expected {
		t.Errorf("Expected output to be '%s', got '%s' instead", expected, result)
	}

	if !performedExit {
		t.Error("Expected app to exit normally but it didn't")
	}
}

func Test_App_HandleError_DoesNothingIfNoError(t *testing.T) {
	var out bytes.Buffer
	log.SetOutput(&out)
	screen := tcell.NewSimulationScreen("")
	app := NewApp(screen, nil, &out)

	app.HandleError(nil, true)
	app.HandleError(nil, false)

	if len(out.String()) > 0 {
		t.Error("Expected no output with a nil error")
	}
}

func Test_App_HandleError_OutputsErrorToLog(t *testing.T) {
	var out, outStd, void bytes.Buffer
	log.SetOutput(&out)
	screen := tcell.NewSimulationScreen("")
	app := NewApp(screen, &void, &outStd)
	app.handleErrorExit = func() {
		// Do nothing for test
	}

	app.HandleError(errors.New("this is an error"), true)

	if len(out.String()) == 0 {
		t.Error("Expected log output with a non-nil error")
	}

	if len(outStd.String()) > 0 {
		t.Error("Expected no output to error stream with logging enabled")
	}
}

func Test_App_HandleError_OutputsErrorToStream(t *testing.T) {
	var out, outLog, void bytes.Buffer
	log.SetOutput(&outLog)
	screen := tcell.NewSimulationScreen("")
	app := NewApp(screen, &void, &out)
	app.handleErrorExit = func() {
		// Do nothing for test
	}

	app.HandleError(errors.New("this is an error"), false)

	if len(out.String()) == 0 {
		t.Error("Expected output with a non-nil error")
	}

	if len(outLog.String()) > 0 {
		t.Error("Expected no log output when logging disabled")
	}
}

func Test_App_HandleError_ExitsScreenBuffer(t *testing.T) {
	var out, void bytes.Buffer
	log.SetOutput(&void)
	screen := tcell.NewSimulationScreen("")
	app := NewApp(screen, &out, &void)
	app.handleErrorExit = func() {
		// Do nothing for test
	}
	app.screenBufferActive = true

	app.HandleError(errors.New("this is an error"), false)
	result := out.String()
	expected := testBufferExitSequence

	if result != expected {
		t.Errorf("Expected hex output to be '%x', got '%x' instead", expected, result)
	}
}

func Test_App_HandleError_ExitsWithError(t *testing.T) {
	var void bytes.Buffer
	log.SetOutput(&void)
	screen := tcell.NewSimulationScreen("")
	app := NewApp(screen, &void, &void)
	exitedOnError := false
	app.handleErrorExit = func() {
		exitedOnError = true
	}

	app.HandleError(errors.New("this is an error"), false)

	if !exitedOnError {
		t.Error("Expected to exit with an error code")
	}
}
