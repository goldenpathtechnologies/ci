package main

import (
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
	"os"
)

type VersionOptions struct {
	Version bool `short:"v" long:"version" description:"Print version information"`
}

type HelpOptions struct {
	Help bool `short:"h" long:"help" description:"Show this help message"`
}

// TODO: Change the name of the AppOptions struct if using the function reference.
//  Options should be data only and not also behaviour.
type AppOptions struct {
	VersionInformation *VersionOptions
	HelpInformation *HelpOptions
	WriteHelp func(writer io.Writer)
}

func GetAppFlags(appName string) (*AppOptions, error) {
	appOptions := &AppOptions{
		VersionInformation: &VersionOptions{},
		HelpInformation: &HelpOptions{},
	}

	parser := flags.NewNamedParser(appName, flags.PrintErrors | flags.PassDoubleDash)

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

	if _, err := parser.AddGroup(
		"Help Options",
		"Help Options",
		appOptions.HelpInformation); err != nil {
		return nil, err
	}

	if _, err := parser.Parse(); err != nil {
		return nil, err
	}

	appOptions.WriteHelp = func(writer io.Writer) {
		parser.WriteHelp(writer)
	}

	return appOptions, nil
}