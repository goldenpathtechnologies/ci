package main

import (
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

type TestOptions struct {
	Test1 bool `short:"a" long:"test1" description:"Test 1"`
	Test2 string `short:"b" long:"test2" description:"Test 2"`
	Test3 string `short:"c" long:"test3" description:"Test 3" value-name:"TEST"`
	Test4 bool `short:"d" long:"test4"`
	Test5 bool `short:"e" long:"test5" description:"This is intended to be a very long description that will wrap to the next line. I'm hoping that this will wrap as many times as necessary to get the effect that I desire. I'm sure that this will only wrap once because there is currently a bug in the help text rendering."`
	Test6 bool `short:"f" long:"test6-really-long-option-that-will-wrap"`
}

type VersionOptions struct {
	Version bool `short:"v" long:"version" description:"Print version information"`
}

func GetAppFlags() (*TestOptions, *VersionOptions, error) {
	testOptions := &TestOptions{}
	versionOptions := &VersionOptions{}

	parser := flags.NewNamedParser("ci", flags.Default)

	parser.UnknownOptionHandler = func(option string, arg flags.SplitArgument, args []string) ([]string, error) {
		parser.WriteHelp(os.Stdout)

		err := flags.Error{
			Type:    flags.ErrUnknownFlag,
			Message: fmt.Sprintf("\nError: unknown flag '%v'", option),
		}

		return nil, errors.New(err.Error())
	}

	if _, err := parser.AddGroup("Test Group", "Long Group", testOptions); err != nil {
		return nil, nil, err
	}

	if _, err := parser.AddGroup("Version Information", "Version Information", versionOptions); err != nil {
		return nil, nil, err
	}

	if _, err := parser.Parse(); err != nil {
		return nil, nil, err
	}

	return testOptions, versionOptions, nil
}