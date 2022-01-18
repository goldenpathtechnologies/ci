package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
)

const maxFilterLength = 32

const (
	filterMethodBeginsWith = iota
	filterMethodEndsWith
	filterMethodContains
	filterMethodGlobPattern
)

const (
	filterTextField = iota
	filterMethodField
)

type FilterForm struct {
	*tview.Form
	filterText   *tview.InputField
	filterMethod *tview.DropDown
	doneHandler  func(key tcell.Key)
}

func CreateFilterForm() *FilterForm {
	filterText := tview.NewInputField().
		SetLabel("Enter filter text:").
		SetFieldWidth(30)

	filterMethod := tview.NewDropDown().
		SetLabel("Filter method:").
		SetOptions(
			[]string{"Begins with", "Ends with", "Contains", "Manual glob"},
			func(text string, index int) {}).
		SetCurrentOption(filterMethodBeginsWith).
		SetListStyles(
			tcell.Style{}.
				Foreground(tcell.ColorBlack).
				Background(tcell.ColorGreen),
			tcell.Style{}.
				Foreground(tcell.ColorBlack).
				Background(tcell.ColorWhite))

	form := tview.NewForm().
		AddFormItem(filterText).
		AddFormItem(filterMethod).
		SetFocus(filterTextField)

	form.SetBorder(true).
		SetTitle("Filter Directory List").
		SetBorderPadding(1, 1, 1, 1)

	filterForm := &FilterForm{
		Form:         form,
		filterText:   filterText,
		filterMethod: filterMethod,
	}

	filterText.SetAcceptanceFunc(filterForm.handleFilterAcceptance)
	form.SetInputCapture(filterForm.handleFilterFormInput)

	return filterForm
}

func (f *FilterForm) handleFilterAcceptance(textToCheck string, lastChar rune) bool {
	if lastChar == '/' || lastChar == '\\' {
		return false
	}

	filterMethod, _ := f.filterMethod.GetCurrentOption()
	globChars := "*?[]!"
	if filterMethod != filterMethodGlobPattern && strings.ContainsRune(globChars, lastChar) {
		return false
	}

	return len(textToCheck) <= maxFilterLength
}

func (f *FilterForm) handleFilterFormInput(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	switch key {
	case tcell.KeyEsc:
		fallthrough
	case tcell.KeyEnter:
		if item, _ := f.GetFocusedItemIndex(); item == filterTextField {
			f.doneHandler(key)
			return nil
		}
		//else if item == filterMethodField {
		//	f.filterMethod.InputHandler()(event, func(p tview.Primitive) {
		// TODO: Collapse dropdown on Enter or Esc and call doneHandler with the key. This
		//  will ensure that the dropdown auto-collapses when the filter is entered. This
		//  will also eliminate the need to TAB back to the filterTextField to enter the filter.
		//  Uncomment this code block and remove this text. Write tests for this as well.
		//	})
		//	f.doneHandler(key)
		//	return nil
		//}
	}

	return event
}

func (f *FilterForm) GetText() string {
	filterText := f.filterText.GetText()
	if filterText == "" {
		return filterText
	}

	currentOptionIndex, _ := f.filterMethod.GetCurrentOption()

	switch currentOptionIndex {
	case filterMethodBeginsWith:
		return filterText + "*"
	case filterMethodEndsWith:
		return "*" + filterText
	case filterMethodContains:
		return "*" + filterText + "*"
	case filterMethodGlobPattern:
		fallthrough
	default:
		return filterText
	}
}

func (f *FilterForm) SetText(text string) *FilterForm {
	f.filterText.SetText(text)

	return f
}

func (f *FilterForm) Clear() {
	f.filterText.SetText("")
}

func (f *FilterForm) SetDoneHandler(handler func(key tcell.Key)) *FilterForm {
	f.doneHandler = handler

	return f
}
