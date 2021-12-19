package ui

import (
	"github.com/rivo/tview"
	"log"
	"os"
)

type App struct {
	*tview.Application
	screenBufferActive bool
	// TODO: Add property to store the output stream (e.g. os.Stdout). This will
	//  make the rest of the application less dependent on a specific output stream.
}

func NewApplication() *App {
	return &App{
		Application:        tview.NewApplication(),
		screenBufferActive: false,
	}
}

// enterScreenBuffer Switches terminal to alternate screen buffer to retain command history
//  of host process
func (a *App) enterScreenBuffer() {
	if a.screenBufferActive {
		return
	}

	print("\033[?1049h")
	a.screenBufferActive = true
}

// exitScreenBuffer Exits the alternate screen buffer and returns to that of host process
func (a *App) exitScreenBuffer() {
	if !a.screenBufferActive {
		return
	}

	print("\033[?1049l")
	a.screenBufferActive = false
}

func (a *App) Start() {
	a.enterScreenBuffer()
}

func (a *App) Stop() {
	a.Application.Stop()
	a.exitScreenBuffer()
}

func (a *App) PrintAndExit(data string) {
	a.Stop()
	// TODO: For testability, use an io.Writer parameter to write output. Looks like
	//  func (a *App) PrintAndExit(output io.Writer, data string). Find other places
	//  in the code that use the output stream instance directly and make similar
	//  changes.
	_, err := os.Stdout.WriteString(data)
	a.HandleError(err, true)
	os.Exit(0)
}

func (a *App) HandleError(err error, logError bool) {
	if err != nil {
		if logError {
			log.Fatal(err)
		}
		a.exitScreenBuffer()
		os.Exit(1)
	}
}
