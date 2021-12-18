package utils

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

type TestApp struct {
	screen tcell.SimulationScreen
	app    *tview.Application
}

func NewTestApp() *TestApp {
	return &TestApp{
		screen: nil,
		app:    nil,
	}
}

func (t *TestApp) Init(width, height int) *TestApp {
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

func (t *TestApp) GetPrimitiveOutput() string {
	x, y, width, height := t.app.GetFocus().GetRect()
	output := ""

	for j := y; j < y + height; j++ {
		for i := x; i < x + width; i++ {
			c, _, _, _ := t.screen.GetContent(i, j)
			output += string(c)
		}
		output += "\n"
	}

	return output
}

func (t *TestApp) Run(p tview.Primitive, autoStop bool, callback func()) {
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

func (t *TestApp) Stop() {
	t.app.Stop()
}
