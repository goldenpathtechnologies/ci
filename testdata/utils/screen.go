package utils

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

type TestScreenApp struct {
	screen tcell.SimulationScreen
	app    *tview.Application
}

func NewTestScreenApp() *TestScreenApp {
	return &TestScreenApp{
		screen: nil,
		app:    nil,
	}
}

func (t *TestScreenApp) Init(width, height int) *TestScreenApp {
	t.screen = setUpSimulationScreen(width, height)
	t.app = setUpTestApplication(t.screen)

	if err := t.screen.Init(); err != nil {
		panic(err)
	}

	return t
}

func setUpSimulationScreen(width, height int) tcell.SimulationScreen {
	simScreen := tcell.NewSimulationScreen("") // "" = UTF-8 charset
	simScreen.SetSize(width, height)
	return simScreen
}

func setUpTestApplication(screen tcell.Screen) *tview.Application {
	return tview.NewApplication().
		SetScreen(screen)
}

func (t *TestScreenApp) GetPrimitiveOutput() string {
	x, y, width, height := t.app.GetFocus().GetRect()
	output := ""

	for j := y; j < y+height; j++ {
		for i := x; i < x+width; i++ {
			c, _, _, _ := t.screen.GetContent(i, j)
			output += string(c)
		}
		output += "\n"
	}

	return output
}

func (t *TestScreenApp) Run(p tview.Primitive, autoStop bool, callback func()) {
	if autoStop {
		defer t.app.Stop()
	}

	go func() {
		if err := t.app.SetRoot(p, false).Run(); err != nil {
			panic(err)
		}
	}()

	// Wait for app to load
	time.Sleep(time.Second)

	callback()
}

func (t *TestScreenApp) Stop() {
	t.app.Stop()
}
