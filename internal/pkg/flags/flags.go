package flags

import (
	"ci/internal/pkg/utils"
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
	"os"
	"time"
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

var (
	BuildVersion string = ""
	BuildDate string = ""
	BuildOwner string = ""
)

func InitFlags(appName string) (*AppOptions, error) {
	var (
		appOptions *AppOptions
		err error
	)

	if appOptions, err = GetAppFlags(appName); err != nil {
		return nil, err
	}

	HandleHelpInformation(appOptions)

	if err = HandleVersionInformation(appOptions, appName); err != nil {
		return nil, err
	}

	return appOptions, nil
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

func HandleHelpInformation(appOptions *AppOptions) {
	if appOptions.HelpInformation.Help {
		appOptions.WriteHelp(os.Stdout)
		os.Exit(0)
	}
}

func HandleVersionInformation(appOptions *AppOptions, appName string) error {
	var (
		buildDate time.Time
		err error
	)

	if appOptions.VersionInformation.Version {
		if buildDate, err = time.Parse(time.RFC3339, BuildDate); err != nil {
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
			appName,
			buildDate.Year(),
			BuildOwner,
			BuildVersion,
			buildDate.Format(time.RFC3339))
		_, err = os.Stdout.WriteString(versionString)
		utils.HandleError(err, true)
		os.Exit(0)
	}

	return nil
}