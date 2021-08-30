package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

const (
	exitCodeErr       = 1
	exitCodeInterrupt = 2
	pathSeparator = string(os.PathSeparator)
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

func GetDirectoryList(location string) ([]string, error) {
	var out bytes.Buffer

	// TODO: Determine which OS is in use and use corresponding commands. Also, consider
	//  using this library instead for cross-platform directory traversal:
	//   https://github.com/karrick/godirwalk
	cmd := exec.Command(
		"powershell.exe",
		"-Command",
		"Get-ChildItem",
		"-Directory",
		fmt.Sprintf(`"%v"`, strings.TrimSpace(location)),
		"|",
		"Select-Object",
		"-ExpandProperty",
		"Name")
	cmd.Stdout = &out
	err := cmd.Run()

	if len(out.String()) == 0 {
		return nil, err
	}

	dirList := strings.Split(strings.TrimSpace(out.String()), "\n")
	for i, _ := range dirList {
		dirList[i] = strings.TrimSpace(dirList[i]) + pathSeparator
	}

	return dirList, err
}

func GetInitialDirectory() (string, error) {
	var out bytes.Buffer

	cmd := exec.Command(
		"powershell.exe",
		"-Command",
		"Get-Location",
		"|",
		"Select-Object",
		"-ExpandProperty",
		"Path")
	cmd.Stdout = &out
	err := cmd.Run()

	dir := strings.TrimSpace(out.String()) + pathSeparator

	return dir, err
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

func GetListUI(app *tview.Application, titleBox *tview.TextView) *tview.List {
	list := tview.NewList().ShowSecondaryText(false)
	tview.Borders.HorizontalFocus = tview.Borders.Horizontal

	list.SetBorder(true)

	currentDir, err := GetInitialDirectory()
	HandleError(err)

	titleBox.Clear()
	titleBox.SetText(currentDir)

	loadList := func(dir string) {
		list.Clear()

		dirs, err := GetDirectoryList(dir)
		HandleError(err)

		if dirs == nil {
			list.AddItem("<Enter directory>", "", 0, nil)
		} else {
			for _, d := range dirs {
				list.AddItem(d,"", 0, nil)
			}
		}


		list.AddItem("<Quit>", "Press to exit", 'q', func() {
			app.Stop()
			ExitScreenBuffer()
		})

		list.AddItem("<Help>", "Get help with this program", 'h', func(){})

		titleBox.Clear()
		titleBox.SetText(currentDir)
	}
	loadList(currentDir)

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			if strings.Count(currentDir, pathSeparator) > 1 {
				log.Println(fmt.Sprintf("Current directory when able to traverse: %v", currentDir))
				paths := strings.Split(currentDir, pathSeparator)
				log.Println(paths)
				paths = paths[:len(paths)-2]
				currentDir = strings.Join(paths, pathSeparator) + pathSeparator
				log.Println(fmt.Sprintf("Current diirectory after mutation: %v", currentDir))
				loadList(currentDir)
			}
			return nil
		case tcell.KeyRight:
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
			log.Fatal(err)
		}
	}
}

func run(app *tview.Application, args []string) error {
	SetBoxBorderStyle()

	titleBox := GetTitleBoxUI()
	list := GetListUI(app, titleBox)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(titleBox, 5, 0, false).
		AddItem(tview.NewFlex().
			AddItem(list, 0, 1, true).
			AddItem(tview.NewBox().
				SetBorder(true).
				SetTitle("Details"),
				0, 1, false),
			0, 1, false)

	if err := app.SetRoot(flex, true).SetFocus(list).Run(); err != nil {
		ExitScreenBuffer()

		return err
	}

	return nil
}

func main() {
	closeLogFile := InitFileLogging()

	defer closeLogFile()

	// Note: code taken from https://pace.dev/blog/2020/02/17/repond-to-ctrl-c-interrupt-signals-gracefully-with-context-in-golang-by-mat-ryer.html
	EnterScreenBuffer()

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
