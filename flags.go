package main

import (
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

type VersionOptions struct {
	Version bool `short:"v" long:"version" description:"Print version information"`
}

type AppOptions struct {
	VersionInformation *VersionOptions
}

func GetAppFlags(appName string) (*AppOptions, error) {
	appOptions := &AppOptions{
		VersionInformation: &VersionOptions{},
	}

	parser := flags.NewNamedParser(appName, flags.Default)

	parser.UnknownOptionHandler = func(option string, arg flags.SplitArgument, args []string) ([]string, error) {
		parser.WriteHelp(os.Stdout)

		err := flags.Error{
			Type:    flags.ErrUnknownFlag,
			Message: fmt.Sprintf("Error: unknown flag '%v'\n", option),
		}

		return nil, errors.New(err.Error())
	}

	if _, err := parser.AddGroup(
		"Version Information",
		"Version Information",
		appOptions.VersionInformation); err != nil {
		return nil, err
	}

	if _, err := parser.Parse(); err != nil {
		return nil, err
	}

	return appOptions, nil
}