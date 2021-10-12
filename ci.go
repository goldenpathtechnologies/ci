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

var (
	BuildVersion string = ""
	BuildDate   string = ""
	testOptions *TestOptions
	versionOptions *VersionOptions
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
		// TODO: Reconsider logging errors as other libraries may also print their errors causing duplication.
		//log.Fatal(err)
		os.Exit(1) // TODO: This is temporary, need to specifically respond to the help option so exit code can be 0.
	}
}

func HandleUIError(err error) {
	if err != nil {
		ExitScreenBuffer()
		// TODO: Reconsider logging errors as other libraries may also print their errors causing duplication.
		//log.Fatal(err)
		os.Exit(1) // TODO: This is temporary, need to specifically respond to the help option so exit code can be 0.
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

func GetDetailsUI() *tview.TextView {
	details := tview.NewTextView()

	details.
		SetDynamicColors(true).
		SetWrap(false).
		SetTitle("Details").
		SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetDrawFunc(GetScrollBarDrawFunc(
			details,
			func() (width, height int) {
				text := details.GetText(true)
				// TODO: This breaks when word wrap is enabled in the TextView as wrapped
				//  lines are not delimited by a new line externally. Ideally, I'd just
				//  access the s.longestLine and s.pageSize fields, but those are
				//  not exported from tview.TextView and therefore inaccessible.
				lines := strings.Split(text, "\n")
				longestLine := ""

				for _, v := range lines {
					if len(v) > len(longestLine) {
						longestLine = v
					}
				}

				return len(longestLine), len(lines)
			},
			func() (vScroll, hScroll int) {
				return details.GetScrollOffset()
			}))

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

func GetListUI(
	app *tview.Application,
	titleBox *tview.TextView,
	filter *tview.InputField,
	pages *tview.Pages,
	details *tview.TextView,
	) *tview.List {

	list := tview.NewList().
		ShowSecondaryText(false).
		SetSelectedTextColor(tcell.ColorBlack)

	list.SetBorder(true).
		SetTitle(listUITitle).
		SetBorderPadding(1, 1, 0, 1).
		SetDrawFunc(GetScrollBarDrawFunc(
			list,
			func() (width, height int) {
				_, _, listWidth, _ := list.GetInnerRect()
				listHeight := list.GetItemCount()

				return listWidth, listHeight
			},
			func() (vScroll, hScroll int) {
				selectedItem := list.GetCurrentItem()
				itemCount := list.GetItemCount()
				_, _, _, pageHeight := list.GetInnerRect()

				v, h := list.GetOffset()

				if selectedItem == 0 {
					v = 0
				}

				if selectedItem == itemCount-1 {
					v = itemCount - pageHeight
				}

				return v, h
			}))

	currentDir, err := GetInitialDirectory()
	HandleUIError(err)

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
			HandleUIError(err)
		})

		scanner, err := godirwalk.NewScanner(currentDir)
		HandleUIError(err)

		for scanner.Scan() {
			d, err := scanner.Dirent()
			HandleUIError(err)

			if d.IsDir() {
				if isMatch, _ := filepath.Match(filterText, d.Name()); len(filterText) == 0 || isMatch {
					list.AddItem(d.Name() + pathSeparator, "", 0, func() {
						path := currentDir + d.Name() + pathSeparator

						app.Stop()
						ExitScreenBuffer()

						_, err := os.Stdout.WriteString(path)
						HandleUIError(err)
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

		getNextItemIndex := func(isIncrementing bool) int {
			var increment int
			if isIncrementing {
				increment = 1
			} else {
				increment = -1
			}
			itemCount := list.GetItemCount()
			nextItemIndex := list.GetCurrentItem() + increment
			// Note: Euclidean modulo operation, https://stackoverflow.com/questions/43018206/modulo-of-negative-integers-in-go
			return ((nextItemIndex % itemCount) + itemCount) % itemCount
		}

		displayDirectoryDetails := func(isNavigatingDown bool) {
			details.Clear()
			nextItemIndex := getNextItemIndex(isNavigatingDown)
			nextItem, _ := list.GetItemText(nextItemIndex)
			if !isMenuItem(nextItem) {
				details.SetText(GetDirectoryInfo(currentDir + nextItem))
			} else if nextItem == listUIEnterDir {
				details.SetText(GetDirectoryInfo(currentDir))
			}
			details.ScrollToBeginning()
		}

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
			displayDirectoryDetails(false)
			return event
		case tcell.KeyDown:
			displayDirectoryDetails(true)
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

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "[red]Unable to read directory details. You may have insufficient privileges.[white]"
	}

	writer := tabwriter.NewWriter(&out, 1, 2, 2, ' ', 0)

	// TODO: Create a function for printing each row of the tab output to reduce duplication.
	_, err = fmt.Fprintf(writer, "%v\t%v\t%v\t%v\n", "Mode", "Name", "ModTime", "Bytes")
	HandleUIError(err)

	_, err = fmt.Fprintf(writer, "%v\t%v\t%v\t%v\n", "----", "----", "-------", "-----")
	HandleUIError(err)

	for _, f := range files {
		dateFormat := "2006-01-02 3:04 PM"
		modTime := f.ModTime().Format(dateFormat)
		_, err := fmt.Fprintf(writer, "%v\t%v\t%v\t%v\n", f.Mode(), f.Name(), modTime, f.Size())
		HandleUIError(err)
	}

	HandleUIError(writer.Flush())

	return out.String()
}

func SetApplicationStyles() {
	// TODO: Setting this interferes with the styles of other components such
	//  as the List. Find a way to target styles to specific components.
	//  Additionally, runes display horribly in PowerShell if not using
	//  Windows Terminal. Find a way to fix this.
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
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
			HandleUIError(err)
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

func InitFlags() {
	var err error

	testOptions, versionOptions, err = GetAppFlags()

	HandleError(err)

	//if testOptions.Help {
	//	//PrintHelpTextAndExit()
	//} else

	if versionOptions.Version {
		_, err = os.Stdout.WriteString("ci version 0.0.0")
		HandleError(err)
		os.Exit(0)
	}
}

func main() {
	InitFlags()

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
		HandleUIError(err)
	}
}
