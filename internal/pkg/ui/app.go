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
	bufferExitSequence = "\033[?1049l"
)

type App struct {
	*tview.Application
	screenBufferActive bool
	outputStream       io.Writer
	errorStream        io.Writer
	handleNormalExit   func()
	handleErrorExit    func()
}

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

func (a *App) PrintAndExit(data string) {
	a.Stop()
	_, err := a.outputStream.Write([]byte(data))
	a.HandleError(err, true)
	a.handleNormalExit()
}

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

// enterScreenBuffer Switches terminal to alternate screen buffer to retain command history
//  of host process
func (a *App) enterScreenBuffer() {
	if a.screenBufferActive {
		return
	}

	if _, err := a.outputStream.Write([]byte(bufferEntrySequence)); err != nil {
		log.Print(err)
		a.handleErrorExit()
	}

	a.screenBufferActive = true
}

// exitScreenBuffer Exits the alternate screen buffer and returns to that of host process
func (a *App) exitScreenBuffer() {
	if !a.screenBufferActive {
		return
	}

	if _, err := a.outputStream.Write([]byte(bufferExitSequence)); err != nil {
		log.Print(err)
		a.handleErrorExit()
	}

	a.screenBufferActive = false
}

func (a *App) Start() {
	a.enterScreenBuffer()
}

func (a *App) Stop() {
	a.Application.Stop()
	a.exitScreenBuffer()
}
