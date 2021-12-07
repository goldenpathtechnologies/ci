package main

import (
	"context"
	"github.com/goldenpathtechnologies/ci/internal/pkg/flags"
	"github.com/goldenpathtechnologies/ci/internal/pkg/ui"
	"github.com/goldenpathtechnologies/ci/internal/pkg/utils"
	"log"
	"os"
	"os/signal"
)

const (
	exitCodeInterrupt = 2
)

var (
	AppName      = "ci"
	BuildVersion string
	BuildDate    string
	BuildOwner   string
)

func main() {
	closeLogFile := utils.InitFileLogging()
	defer closeLogFile()

	var (
		options = &flags.AppOptions{
			AppName:      AppName,
			BuildVersion: BuildVersion,
			BuildDate:    BuildDate,
			BuildOwner:   BuildOwner,
		}
		err error
	)

	if options, err = flags.InitFlags(options); err != nil {
		fErr, isFlagError := err.(*flags.FlagError)
		if isFlagError {
			if fErr.ErrorCode == flags.FlagErrorNormalExit {
				os.Exit(0)
			} else {
				log.Fatal(fErr)
			}
		} else {
			log.Fatal(err)
		}
	}

	app := ui.NewApplication()
	app.Start()


	// Note: code taken and modified from https://pace.dev/blog/2020/02/17/repond-to-ctrl-c-interrupt-signals-gracefully-with-context-in-golang-by-mat-ryer.html
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	defer func() {
		signal.Stop(signalChan)
		cancel()
		app.Stop()
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

	if err = ui.Run(app, options); err != nil {
		app.HandleError(err, true)
	}
}