package ui

import (
	"github.com/goldenpathtechnologies/ci/internal/pkg/options"
	"regexp"
	"testing"
	"time"
)

func Test_Help_GetHelpText_BuildDateDefaultsToNowWhenInputFormatIsInvalid(t *testing.T) {
	appOptions := &options.AppOptions{}

	// Note: The Time values we're working with have a low enough precision to be
	//  sufficiently deterministic. Abstract all direct references to Time as soon
	//  as this proves false.
	currentDate := time.Now().Format(time.RFC3339)
	helpText := GetHelpText(appOptions)
	match := regexp.MustCompile("Build date: (.*)")
	matches := match.FindStringSubmatch(helpText)
	helpTextBuildDate := matches[1]

	if currentDate != helpTextBuildDate {
		t.Errorf("Expected the build date to be '%s', got '%s' instead", currentDate, helpTextBuildDate)
	}
}

func Test_Help_GetHelpText_ReturnsEmptyStringWhenAppOptionsIsNil(t *testing.T) {
	result := GetHelpText(nil)

	if result != "" {
		t.Errorf("Expected an empty string but got '%s' instead", result)
	}
}