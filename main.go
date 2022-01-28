package main

import (
	"context"
	"github.com/goldenpathtechnologies/ci/internal/pkg/options"
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
	BuildOwner1  string
	BuildOwner2  string
	Repository   = "https://github.com/goldenpathtechnologies/ci"
)

func main() {
	closeLogFile := utils.InitFileLogging()
	defer closeLogFile()

	var (
		appOptions = &options.AppOptions{
			AppName:      AppName,
			BuildVersion: BuildVersion,
			BuildDate:    BuildDate,
			BuildOwner1:  BuildOwner1,
			BuildOwner2:  BuildOwner2,
			Repository:   Repository,
		}
		err error
	)

	if appOptions, err = appOptions.Init(); err != nil {
		fErr, isFlagError := err.(*options.OptionError)
		if isFlagError {
			if fErr.ErrorCode == options.OptionErrorNormalExit {
				os.Exit(0)
			} else {
				log.Fatal(fErr)
			}
		} else {
			log.Fatal(err)
		}
	}

	app := ui.NewApp(nil, os.Stdout, os.Stderr)
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

	if err = ui.Run(app, appOptions); err != nil {
		app.HandleError(err, true)
	}
}
