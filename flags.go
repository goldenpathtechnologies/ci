package main

import (
	"bytes"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/go-wordwrap"
	"reflect"
	"strings"
	"text/tabwriter"
)

type FlagOptions struct {
	Version bool `short:"v" long:"version" description:"Print version information"`
	//Help bool `short:"h" long:"help" description:"Print command line options"`
}

func GetAppFlags() (*FlagOptions, error) {
	options := &FlagOptions{}

	_, err := flags.Parse(options)

	return options, err
}

func GetHelpText() (string, error) {
	var out bytes.Buffer

	writer := tabwriter.NewWriter(&out, 0, 2, 2, ' ', 0)

	t := reflect.TypeOf(FlagOptions{})

	for _, fieldName := range []string{"Version", "Help"} {
		field, found := t.FieldByName(fieldName)

		if !found {
			continue
		}

		shortName := field.Tag.Get("short")
		longName := field.Tag.Get("long")
		description := field.Tag.Get("description")
		valueName := field.Tag.Get("value-name")

		if field.Type.Kind() != reflect.Bool {
			if len(valueName) > 0 {
				shortName = fmt.Sprintf("-%v=%v", shortName, valueName)
				longName = fmt.Sprintf("-%v=%v", longName, valueName)
			} else {
				shortName = fmt.Sprintf("-%v=...", shortName)
				longName = fmt.Sprintf("-%v=...", longName)
			}
		}

		maxOptionColWidth := 40
		maxDescriptionColWidth := 60

		optionName := fmt.Sprintf("%v, %v", shortName, longName)
		optionName = wordwrap.WrapString(optionName, uint(maxOptionColWidth))
		optionLines := strings.Split(optionName, "\n")

		description = wordwrap.WrapString(description, uint(maxDescriptionColWidth))
		descriptionLines := strings.Split(description, "\n")

		if len(optionLines) > 0 && len(descriptionLines) > 0 {
			_, err := fmt.Fprintf(writer, "%v\t%v", optionLines[0], descriptionLines[0])

			if err != nil {
				return "", err
			}

			if len(optionLines) > 1 || len(descriptionLines) > 1 {
				var o, d string

				if len(optionLines) > 1 {
					o = optionLines[1]
				} else {
					o = ""
				}

				if len(descriptionLines) > 1 {
					d = descriptionLines[1]
				} else {
					d = ""
				}

				_, err := fmt.Fprintf(writer, "%v\t%v", o, d)

				if err != nil {
					return "", err
				}
			}
		}
	}

	err := writer.Flush()

	if err != nil {
		return "", err
	}

	helpTextHeader := `Usage: ci [options]

Options:
`
	return helpTextHeader + out.String(), nil
}