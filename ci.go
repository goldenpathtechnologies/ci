package main

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/karrick/godirwalk"
	"github.com/rivo/tview"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
)

const (
	exitCodeErr       = 1
	exitCodeInterrupt = 2
	pathSeparator = string(os.PathSeparator)
	listUITitle = "Directory List"
)

func HandleError(err error) {
	if err != nil {
		ExitScreenBuffer()
		log.Fatal(err)
	}
}

// EnterScreenBuffer Switches terminal to alternate screen buffer to retain command history
//  of host process
func EnterScreenBuffer() {
	print("\033[?1049h")
}

// ExitScreenBuffer Exits the alternate screen buffer and returns to that of host process
func ExitScreenBuffer() {
	print("\033[?1049l")
}

func GetInitialDirectory() (string, error) {
	dir, err := filepath.Abs(".")

	return dir + pathSeparator, err
}

func GetFilterUI() *tview.InputField {
	// TODO: Put the field width and max length in a constant
	filter := tview.NewInputField().
		SetLabel("Enter filter text: ").
		SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
			if lastChar == '/' || lastChar == '\\' {
				return false
			}

			return len(textToCheck) <= 32
		})

	filter.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetTitle("Filter Directory List")

	return filter
}

func CreateModalUI(widget tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(widget, 1, 1, 1, 1, 0, 0, false)
}

func GetTitleBoxUI() *tview.TextView {
	titleBox := tview.NewTextView().
		SetText("No directory has been selected yet!").
		SetTextAlign(tview.AlignCenter).
		SetScrollable(false)

	titleBox.SetBorder(true).
		SetTitle(`CI - "CD Improved"`).
		SetBorderPadding(1,1,1,1).
		SetTitleColor(tcell.ColorGreen)

	return titleBox
}

func GetListUI(app *tview.Application, titleBox *tview.TextView, filter *tview.InputField, pages *tview.Pages) *tview.List {
	list := tview.NewList().ShowSecondaryText(false)
	tview.Borders.HorizontalFocus = tview.Borders.Horizontal

	list.SetBorder(true).SetTitle(listUITitle)

	currentDir, err := GetInitialDirectory()
	HandleError(err)

	titleBox.Clear()
	titleBox.SetText(currentDir)

	var filterText string

	loadList := func(dir string) {
		list.Clear()

		scanner, err := godirwalk.NewScanner(currentDir)
		HandleError(err)

		for scanner.Scan() {
			d, err := scanner.Dirent()
			HandleError(err)

			if d.IsDir() {
				if isMatch, _ := filepath.Match(filterText, d.Name()); len(filterText) == 0 || isMatch {
					list.AddItem(d.Name() + pathSeparator, "", 0, func() {
						path := currentDir + d.Name() + pathSeparator

						app.Stop()
						ExitScreenBuffer()

						_, err := os.Stdout.WriteString(path)
						HandleError(err)
					})
				}
			}
		}

		if list.GetItemCount() == 0 {
			list.AddItem("<Enter directory>", "", 0, func() {
				app.Stop()
				ExitScreenBuffer()

				_, err := os.Stdout.WriteString(currentDir)
				HandleError(err)
			})
		}

		list.AddItem("<Quit>", "Press to exit", 'q', func() {
			app.Stop()
			ExitScreenBuffer()
		})

		list.AddItem("<Help>", "Get help with this program", 'h', func(){})

		list.AddItem("<Filter>", "Filter directories by text", 'f', func() {
			pages.ShowPage("Filter")
			app.SetFocus(filter)
		})

		titleBox.Clear()
		titleBox.SetText(currentDir)
	}
	loadList(currentDir)

	filter.SetDoneFunc(func(key tcell.Key) {
		filterText = filter.GetText()
		if len(filterText) > 0 {
			list.SetTitle(fmt.Sprintf("%v - Filter: %v", listUITitle, filterText))
		} else {
			list.SetTitle(listUITitle)
		}
		filter.SetText("")
		pages.HidePage("Filter")
		app.SetFocus(list)
		loadList(currentDir)
	})

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			if strings.Count(currentDir, pathSeparator) > 1 {
				filterText = ""
				paths := strings.Split(currentDir, pathSeparator)
				paths = paths[:len(paths)-2]
				currentDir = strings.Join(paths, pathSeparator) + pathSeparator
				loadList(currentDir)
			}
			return nil
		case tcell.KeyRight:
			filterText = ""
			selectedDir, _ := list.GetItemText(list.GetCurrentItem())
			currentDir = currentDir +selectedDir
			loadList(currentDir)
			return nil
		}

		return event
	})

	return list
}

func SetBoxBorderStyle() {
	tview.Borders.HorizontalFocus = tview.Borders.Horizontal
	tview.Borders.VerticalFocus = tview.Borders.Vertical
	tview.Borders.TopLeftFocus = tview.Borders.TopLeft
	tview.Borders.TopRightFocus = tview.Borders.TopRight
	tview.Borders.BottomLeftFocus = tview.Borders.BottomLeft
	tview.Borders.BottomRightFocus = tview.Borders.BottomRight
}

// InitFileLogging Initializes logging to a file and returns the function that closes that file
func InitFileLogging() func() {
	file, err := os.OpenFile("./.log", os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	return func() {
		if err := file.Close(); err != nil {
			log.SetOutput(os.Stdout)
			HandleError(err)
		}
	}
}

func run(app *tview.Application, args []string) error {
	SetBoxBorderStyle()

	pages := tview.NewPages()

	filter := GetFilterUI()
	titleBox := GetTitleBoxUI()
	list := GetListUI(app, titleBox, filter, pages)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(titleBox, 5, 0, false).
		AddItem(tview.NewFlex().
			AddItem(list, 0, 1, true).
			AddItem(tview.NewBox().
				SetBorder(true).
				SetTitle("Details"),
				0, 1, false),
			0, 1, false)

	pages.AddPage("Home", flex, true, true).
		AddPage("Filter", CreateModalUI(filter, 40, 7), true, false)

	if err := app.SetRoot(pages, true).SetFocus(list).Run(); err != nil {
		ExitScreenBuffer()

		return err
	}

	return nil
}

func main() {
	closeLogFile := InitFileLogging()
	defer closeLogFile()

	EnterScreenBuffer()

	// Note: code taken from https://pace.dev/blog/2020/02/17/repond-to-ctrl-c-interrupt-signals-gracefully-with-context-in-golang-by-mat-ryer.html
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	app := tview.NewApplication()

	defer func() {
		signal.Stop(signalChan)
		cancel()
		app.Stop()
		ExitScreenBuffer()
	}()

	go func() {
		select {
		case <-signalChan: // first signal, cancel context
			cancel()
		case <-ctx.Done():
		}
		<-signalChan // second signal, hard exit
		os.Exit(exitCodeInterrupt)
	}()
	if err := run(app, os.Args); err != nil {
		HandleError(err)
	}
}
