package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/goldenpathtechnologies/ci/internal/pkg/options"
	"time"
)

// GetHelpText returns the text of the in-app help info.
func GetHelpText(options *options.AppOptions) string {
	var (
		copyrightYear int
		buildDate     time.Time
		err           error
	)

	if options == nil {
		return ""
	}

	if buildDate, err = time.Parse(time.RFC3339, options.BuildDate); err != nil {
		buildDate = time.Now()
	}

	copyrightYear = buildDate.Year()

	helpText := fmt.Sprintf(`[yellow]Directory List[white]
[green]%c[white]        Navigate to child directory (when selected)
[green]%c[white]        Navigate to parent directory
[green]%c[white]        Select previous item
[green]%c[white]        Select next item
[green]%s[white]     Select first item on previous page
[green]%s[white]     Select last item on next page
[green]%s[white]    Enter selected directory/Select option
[green]%s[white]      Select the details pane
[green]%s[white]        Exit and navigate to current directory
[green]%s[white]        Open filter dialog
[green]%s[white]        Show this help text
[green]%s[white]        Exit without navigating

[yellow]Details/Help[white]
[green]%s[white]   Scroll text
[green]%s[white]     Scroll to previous page
[green]%s[white]     Scroll to next page
[green]%s[white]    Deselect back to the directory list
[green]%s[white]
[green]%s[white]

[yellow]Filter[white]
[green]%s[white]    Enter filter text, or clear the existing filter if empty
[green]%s[white]      Set focus to next input field
[green]%s[white]    Expand/collapse filter method dropdown when selected
[green]%s[white]  Select a filter method when dropdown is expanded


[green]%s[white]
Version: %s
Build date: %s
Repository: %s

Copyright © %v
%s
%s


This program is MIT licensed.`,
		tcell.RuneRArrow,
		tcell.RuneLArrow,
		tcell.RuneUArrow,
		tcell.RuneDArrow,
		"PgUp",
		"PgDn",
		"ENTER",
		"TAB",
		"e",
		"f",
		"h",
		"q",
		"ARROWS",
		"PgUp",
		"PgDn",
		"ENTER",
		"TAB",
		"ESC",
		"ENTER",
		"TAB",
		"SPACE",
		"UP/DOWN",
		options.AppName,
		options.BuildVersion,
		buildDate.Format(time.RFC3339),
		options.Repository,
		copyrightYear,
		options.BuildOwner1,
		options.BuildOwner2)

	return helpText
}


