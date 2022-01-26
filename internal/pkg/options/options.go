// Package options implements command line options and manages application wide data.
package options

import (
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"time"
)

// VersionOptions defines the properties of the '-v' or '--version' command line option.
type VersionOptions struct {
	Version bool `short:"v" long:"version" description:"Print version information"`
}

// HelpOptions defines the properties of the '-h' or '--help' command line option.
type HelpOptions struct {
	Help bool `short:"h" long:"help" description:"Show this help message"`
}

// AppOptions stores information that is used throughout the application.
type AppOptions struct {
	VersionInformation *VersionOptions
	HelpInformation    *HelpOptions
	AppName            string
	BuildVersion       string
	BuildDate          string
	BuildOwner1        string
	BuildOwner2        string
	Repository         string
}

const (
	OptionErrorNormalExit = iota
	OptionErrorUnexpected
)

// OptionError represents an error that occurred while parsing or handling command line options.
type OptionError struct {
	Err error
	ErrorCode int
}

// Error returns the error message in OptionError.
func (f *OptionError) Error() string {
	if f.ErrorCode == OptionErrorUnexpected {
		return f.Err.Error()
	}

	return ""
}

// Init parses command line options and updates the AppOptions data. Depending on the option,
// Init either handles it (e.g. --help and --version) or sets its data.
func (a *AppOptions) Init() (*AppOptions, error) {
	a.VersionInformation = &VersionOptions{}
	a.HelpInformation = &HelpOptions{}

	parser := flags.NewNamedParser(a.AppName, flags.PrintErrors | flags.PassDoubleDash)

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
		a.VersionInformation); err != nil {
		return nil, err
	}

	if _, err := parser.AddGroup(
		"Help Options",
		"Help Options",
		a.HelpInformation); err != nil {
		return nil, err
	}

	if _, err := parser.Parse(); err != nil {
		return nil, err
	}

	if err := a.handleHelpInformation(parser); err != nil {
		return nil, err
	}

	if err := a.handleVersionInformation(); err != nil {
		return nil, err
	}

	return a, nil
}

// handleHelpInformation prints help information and returns a graceful exit error.
func (a *AppOptions) handleHelpInformation(parser *flags.Parser) error {
	if a.HelpInformation.Help {
		parser.WriteHelp(os.Stdout)
		return &OptionError{ErrorCode: OptionErrorNormalExit}
	}

	return nil
}

// handleVersionInformation prints version information and returns a graceful exit error.
func (a *AppOptions) handleVersionInformation() error {
	var (
		buildDate time.Time
		err error
	)

	if a.VersionInformation.Version {
		if buildDate, err = time.Parse(time.RFC3339, a.BuildDate); err != nil {
			return err
		}

		versionFormat := `%s
Version: %s
Build date: %s
Repository: %s

Copyright © %v
%s
%s
`
		versionString := fmt.Sprintf(
			versionFormat,
			a.AppName,
			a.BuildVersion,
			buildDate.Format(time.RFC3339),
			a.Repository,
			buildDate.Year(),
			a.BuildOwner1,
			a.BuildOwner2)

		if _, err = os.Stdout.WriteString(versionString); err != nil {
			log.Fatal(err)
		}

		return &OptionError{ErrorCode: OptionErrorNormalExit}
	}

	return nil
}
