package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/karrick/godirwalk"
	"github.com/rivo/tview"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

const (
	exitCodeErr       = 1
	exitCodeInterrupt = 2
	pathSeparator = string(os.PathSeparator)
	listUITitle = "Directory List"
	listUIQuit = "<Quit>"
	listUIHelp = "<Help>"
	listUIFilter = "<Filter>"
	listUIEnterDir = "<Enter directory>"
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

func GetDetailsUI() *ScrollView {
	details := NewScrollView()

	details.
		SetBorder(true).
		SetDynamicColors(true).
		SetWrap(false).
		SetTitle("Details").
		SetBorderPadding(1, 1, 1, 1)

	return details
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

func GetListUI(app *tview.Application, titleBox *tview.TextView, filter *tview.InputField, pages *tview.Pages, details *ScrollView, ) *tview.List {

	list := tview.NewList().ShowSecondaryText(false)

	list.SetBorder(true).
		SetTitle(listUITitle).
		SetBorderPadding(1, 1, 0, 1)

	currentDir, err := GetInitialDirectory()
	HandleError(err)

	titleBox.Clear()
	titleBox.SetText(currentDir)

	details.Clear().
		SetText(GetDirectoryInfo(currentDir)).
		ScrollToBeginning().
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyEscape:
				fallthrough
			case tcell.KeyEnter:
				fallthrough
			case tcell.KeyTab:
				app.SetFocus(list)
				return nil
			case 'q':
				app.Stop()
				ExitScreenBuffer()
			}

			return event
		})

	var filterText string

	menuItems := map[string]string{
		listUIQuit: listUIQuit,
		listUIHelp: listUIHelp,
		listUIFilter: listUIFilter,
		listUIEnterDir: listUIEnterDir,
	}

	isMenuItem := func(text string) bool {
		_, exists := menuItems[text]
		return exists
	}

	loadList := func(dir string) {
		list.Clear()

		list.AddItem(listUIEnterDir, "", 'e', func() {
			app.Stop()
			ExitScreenBuffer()

			_, err := os.Stdout.WriteString(currentDir)
			HandleError(err)
		})

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

		list.AddItem(listUIQuit, "Press to exit", 'q', func() {
			app.Stop()
			ExitScreenBuffer()
		})

		list.AddItem(listUIHelp, "Get help with this program", 'h', func(){})

		list.AddItem(listUIFilter, "Filter directories by text", 'f', func() {
			pages.ShowPage("Filter")
			app.SetFocus(filter)
		})

		titleBox.Clear()
		titleBox.SetText(currentDir)
	}
	loadList(currentDir)

	// TODO: Ensure Esc does not apply any existing filter
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
		selectedItem, _ := list.GetItemText(list.GetCurrentItem())

		switch event.Key() {
		case tcell.KeyLeft:
			if strings.Count(currentDir, pathSeparator) > 1 {
				filterText = ""
				list.SetTitle(listUITitle)
				paths := strings.Split(currentDir, pathSeparator)
				paths = paths[:len(paths)-2]
				currentDir = strings.Join(paths, pathSeparator) + pathSeparator
				loadList(currentDir)

				details.Clear()
				details.SetText(GetDirectoryInfo(currentDir)).
					ScrollToBeginning()
			}
			return nil
		case tcell.KeyRight:
			if !isMenuItem(selectedItem) {
				filterText = ""
				list.SetTitle(listUITitle)
				nextDir := currentDir + selectedItem
				if DirectoryIsAccessible(nextDir) {
					currentDir = nextDir
					loadList(currentDir)
				} else {
					details.Clear().
						SetText("[red]Directory inaccessible, unable to navigate. You may have insufficient privileges.[white]").
						ScrollToBeginning()
				}
			}
			return nil
		case tcell.KeyUp:
			details.Clear()
			itemCount := list.GetItemCount()
			nextItemIndex := list.GetCurrentItem() - 1
			// Note: Euclidean modulo operation, https://stackoverflow.com/questions/43018206/modulo-of-negative-integers-in-go
			nextItemIndex = ((nextItemIndex % itemCount) + itemCount) % itemCount
			nextItem, _ := list.GetItemText(nextItemIndex)
			if !isMenuItem(nextItem) {
				details.SetText(GetDirectoryInfo(currentDir + nextItem))
			} else if nextItem == listUIEnterDir {
				details.SetText(GetDirectoryInfo(currentDir))
			}
			details.ScrollToBeginning()
			return event
		case tcell.KeyDown:
			// TODO: This code is somewhat duplicated from the KeyUp case. Use Euclidean modulus operation here as well.
			details.Clear()
			nextItem, _ := list.GetItemText((list.GetCurrentItem() + 1) % list.GetItemCount())
			if !isMenuItem(nextItem) {
				details.SetText(GetDirectoryInfo(currentDir + nextItem))
			} else if nextItem == listUIEnterDir {
				details.SetText(GetDirectoryInfo(currentDir))
			}
			details.ScrollToBeginning()
			return event
		case tcell.KeyTab:
			app.SetFocus(details)
			return nil
		}

		return event
	})

	return list
}

func DirectoryIsAccessible(dir string) bool {
	_, err := ioutil.ReadDir(dir)

	return err == nil
}

func GetDirectoryInfo(dir string) string {
	var out bytes.Buffer

	// TODO: Return an empty string or an error message for directories that have elevated permissions.
	//  Just don't quit the program when this occurs since it is poor UX.
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "[red]Unable to read directory details. You may have insufficient privileges.[white]"
	}

	writer := tabwriter.NewWriter(&out, 1, 2, 2, ' ', 0)

	// TODO: Create a function for printing each row of the tab output to reduce duplication.
	_, err = fmt.Fprintf(writer, "%v\t%v\t%v\t%v\n", "Mode", "Name", "ModTime", "Bytes")
	HandleError(err)

	_, err = fmt.Fprintf(writer, "%v\t%v\t%v\t%v\n", "----", "----", "-------", "-----")
	HandleError(err)

	for _, f := range files {
		dateFormat := "2006-01-02 3:04 PM"
		modTime := f.ModTime().Format(dateFormat)
		_, err := fmt.Fprintf(writer, "%v\t%v\t%v\t%v\n", f.Mode(), f.Name(), modTime, f.Size())
		HandleError(err)
	}

	HandleError(writer.Flush())

	return out.String()
}

func SetApplicationStyles() {
	//tview.Borders.HorizontalFocus = tview.Borders.Horizontal
	//tview.Borders.VerticalFocus = tview.Borders.Vertical
	//tview.Borders.TopLeftFocus = tview.Borders.TopLeft
	//tview.Borders.TopRightFocus = tview.Borders.TopRight
	//tview.Borders.BottomLeftFocus = tview.Borders.BottomLeft
	//tview.Borders.BottomRightFocus = tview.Borders.BottomRight

	// TODO: Setting this interferes with the styles of other components such
	//  as the List. Find a way to target styles to specific components.
	//  Additionally, runes display horribly in PowerShell if not using
	//  Windows Terminal. Find a way to fix this.
	//tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault

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
	SetApplicationStyles()

	pages := tview.NewPages()

	filter := GetFilterUI()
	details := GetDetailsUI()
	titleBox := GetTitleBoxUI()
	list := GetListUI(app, titleBox, filter, pages, details)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(titleBox, 5, 0, false).
		AddItem(tview.NewFlex().
			AddItem(list, 0, 1, true).
			AddItem(details, 0, 2, false),
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
