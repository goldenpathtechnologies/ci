package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"io"
	"log"
	"os"
)

const (
	bufferEntrySequence = "\033[?1049h"
	bufferExitSequence  = "\033[?1049l"
)

// App is an abstraction of the tview.Application with additional functionality.
type App struct {
	*tview.Application
	screenBufferActive bool
	outputStream       io.Writer
	errorStream        io.Writer
	handleNormalExit   func()
	handleErrorExit    func()
	// TODO: Add a flag that enables/disables logging throughout the app so that
	//  it is handled consistently. I discovered during testing that I have to assume
	//  the SUT enabled logging to determine where error output is received. It would
	//  be better to configure this during tests so that I know where error output will
	//  be at any time. See Test_DirectoryList_getDetailsText_HandlesUnexpectedErrors
	//  for issues related to this change.
}

// NewApp creates a new instance of an App.
func NewApp(screen tcell.Screen, stream io.Writer, errStream io.Writer) *App {
	return &App{
		Application:        tview.NewApplication().SetScreen(screen),
		screenBufferActive: false,
		outputStream:       stream,
		errorStream:        errStream,
		handleNormalExit: func() {
			os.Exit(0)
		},
		handleErrorExit: func() {
			os.Exit(1)
		},
	}
}

// PrintAndExit prints data to the App's configured output stream and exits the program.
func (a *App) PrintAndExit(data string) {
	a.Stop()
	_, err := a.outputStream.Write([]byte(data))
	a.HandleError(err, true)
	a.handleNormalExit()
}

// HandleError logs errors and gracefully exits the program with a code of 1.
func (a *App) HandleError(err error, logError bool) {
	if err != nil {
		if logError {
			log.Print(err)
			a.handleErrorExit()
		} else {
			if _, pErr := a.errorStream.Write([]byte(err.Error())); pErr != nil {
				panic(pErr)
			}
		}
		a.exitScreenBuffer()
		a.handleErrorExit()
	}
}

// enterScreenBuffer switches terminal to alternate screen buffer to retain command history
//  of host process
func (a *App) enterScreenBuffer() {
	if a.screenBufferActive {
		return
	}

	if _, err := a.errorStream.Write([]byte(bufferEntrySequence)); err != nil {
		log.Print(err)
		a.handleErrorExit()
	}

	a.screenBufferActive = true
}

// exitScreenBuffer exits the alternate screen buffer and returns to that of host process.
func (a *App) exitScreenBuffer() {
	if !a.screenBufferActive {
		return
	}

	if _, err := a.errorStream.Write([]byte(bufferExitSequence)); err != nil {
		log.Print(err)
		a.handleErrorExit()
	}

	a.screenBufferActive = false
}

// Start starts the application by switching to an alternate screen buffer.
func (a *App) Start() {
	a.enterScreenBuffer()
}

// Stop stops the application and switches to the screen buffer of the host process.
func (a *App) Stop() {
	a.Application.Stop()
	a.exitScreenBuffer()
}
