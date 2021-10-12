package main

import (
	"bytes"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/go-wordwrap"
	"log"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"
)

type FlagOptions struct {
	Version bool `short:"v" long:"version" description:"Print version information"`
	Help bool `short:"h" long:"help" description:"Print command line options"`
	Test1 bool `short:"a" long:"test1" description:"Test 1"`
	Test2 string `short:"b" long:"test2" description:"Test 2"`
	Test3 string `short:"c" long:"test3" description:"Test 3" value-name:"TEST"`
	Test4 bool `short:"d" long:"test4"`
}

func GetAppFlags() (*FlagOptions, error) {
	options := &FlagOptions{}

	// TODO: Reconsider enabling error printing as I may want to handle them myself.
	parser := flags.NewParser(options, flags.PrintErrors | flags.PassDoubleDash)

	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrUnknownFlag {
				PrintHelpTextAndExit()
			}
		default:
			return nil, err
		}
	}

	return options, nil
}

func PrintHelpTextAndExit() {
	if helpText, err := GetHelpText(); err != nil {
		log.Fatal(err)
	} else if _, err = os.Stdout.WriteString(helpText); err != nil {
		log.Fatal(err)
	} else {
		os.Exit(0)
	}
}

func GetHelpText() (string, error) {
	var out bytes.Buffer

	writer := tabwriter.NewWriter(&out, 0, 2, 2, ' ', 0)

	t := reflect.TypeOf(FlagOptions{})

	// TODO: Use the poorly documented 'flags' library functions to iterate on options instead of using reflect.
	for _, fieldName := range []string{"Version", "Help", "Test1", "Test2", "Test3", "Test4"} {
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
				shortName = fmt.Sprintf("-%v %v", shortName, valueName)
				longName = fmt.Sprintf("--%v=%v", longName, valueName)
			} else {
				shortName = fmt.Sprintf("-%v", shortName)
				longName = fmt.Sprintf("--%v=...", longName)
			}
		} else {
			shortName = fmt.Sprintf("-%v", shortName)
			longName = fmt.Sprintf("--%v", longName)
		}

		maxOptionColWidth := 40
		maxDescriptionColWidth := 60

		optionName := fmt.Sprintf("%v, %v", shortName, longName)
		optionName = wordwrap.WrapString(optionName, uint(maxOptionColWidth))
		optionLines := strings.Split(optionName, "\n")

		description = wordwrap.WrapString(description, uint(maxDescriptionColWidth))
		descriptionLines := strings.Split(description, "\n")

		if len(optionLines) > 0 && len(descriptionLines) > 0 {
			_, err := fmt.Fprintf(writer, "  %v\t%v\n", optionLines[0], descriptionLines[0])

			if err != nil {
				return "", err
			}

			// TODO: Need to iterate on the number of additional lines max(optionLines, descriptionLines) times.
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

				_, err := fmt.Fprintf(writer, "  %v\t%v\n", o, d)

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

	// TODO: Make this application agnostic by removing references to specific programs.
	helpTextHeader := `Usage: ci [options]

Options:
`
	return helpTextHeader + out.String(), nil
}