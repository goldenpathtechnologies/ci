package flags

import (
	"ci/internal/pkg/utils"
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"time"
)

type VersionOptions struct {
	Version bool `short:"v" long:"version" description:"Print version information"`
}

type HelpOptions struct {
	Help bool `short:"h" long:"help" description:"Show this help message"`
}

type AppOptions struct {
	VersionInformation *VersionOptions
	HelpInformation    *HelpOptions
	AppName            string
	BuildVersion       string
	BuildDate          string
	BuildOwner         string
}

const (
	FlagErrorUnexpected = iota
	FlagErrorNormalExit
)

type FlagError struct {
	Err error
	ErrorCode int
}

func (f *FlagError) Error() string {
	if f.ErrorCode == FlagErrorUnexpected {
		return f.Err.Error()
	}

	return ""
}

func InitFlags(options *AppOptions) (*AppOptions, error) {
	var err error

	if options, err = getAppFlags(options); err != nil {
		return nil, err
	}

	return options, nil
}

func getAppFlags(options *AppOptions) (*AppOptions, error) {
	options.VersionInformation = &VersionOptions{}
	options.HelpInformation = &HelpOptions{}

	parser := flags.NewNamedParser(options.AppName, flags.PrintErrors | flags.PassDoubleDash)

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
		options.VersionInformation); err != nil {
		return nil, err
	}

	if _, err := parser.AddGroup(
		"Help Options",
		"Help Options",
		options.HelpInformation); err != nil {
		return nil, err
	}

	if _, err := parser.Parse(); err != nil {
		return nil, err
	}

	if err := handleHelpInformation(options, parser); err != nil {
		return nil, err
	}

	if err := handleVersionInformation(options); err != nil {
		return nil, err
	}

	return options, nil
}

func handleHelpInformation(options *AppOptions, parser *flags.Parser) error {
	if options.HelpInformation.Help {
		parser.WriteHelp(os.Stdout)
		return &FlagError{ErrorCode: FlagErrorNormalExit}
	}

	return nil
}

func handleVersionInformation(options *AppOptions) error {
	var (
		buildDate time.Time
		err error
	)

	if options.VersionInformation.Version {
		if buildDate, err = time.Parse(time.RFC3339, options.BuildDate); err != nil {
			return err
		}

		versionFormat := `%v
Copyright © %v
%v

Version: %v
Build date: %v
`
		versionString := fmt.Sprintf(
			versionFormat,
			options.AppName,
			buildDate.Year(),
			options.BuildOwner,
			options.BuildVersion,
			buildDate.Format(time.RFC3339))
		_, err = os.Stdout.WriteString(versionString)
		utils.HandleError(err, true)
		return &FlagError{ErrorCode: FlagErrorNormalExit}
	}

	return nil
}
