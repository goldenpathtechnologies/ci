package main

import (
	"ci/internal/pkg/utils"
	"context"
	"github.com/rivo/tview"
	"os"
	"os/signal"
)

const (
	exitCodeErr       = 1
	exitCodeInterrupt = 2
)

func main() {
	InitFlags()

	closeLogFile := utils.InitFileLogging()
	defer closeLogFile()

	utils.EnterScreenBuffer()

	// Note: code taken and modified from https://pace.dev/blog/2020/02/17/repond-to-ctrl-c-interrupt-signals-gracefully-with-context-in-golang-by-mat-ryer.html
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	app := tview.NewApplication()

	defer func() {
		signal.Stop(signalChan)
		cancel()
		app.Stop()
		utils.ExitScreenBuffer()
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
		utils.HandleError(err, true)
	}
}