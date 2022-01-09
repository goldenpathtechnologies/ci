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
		SetLabel("Enter filter text: ").
		SetFieldWidth(30)

	filterMethod := tview.NewDropDown().
		SetLabel("Filter method").
		SetOptions(
			[]string{"Begins with", "Ends with", "Contains", "Manual glob"},
			func(text string, index int) {}).
		SetCurrentOption(filterMethodBeginsWith)

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

	// TODO: Prevent glob characters from being accepted unless the manual glob filter
	//  method is selected. This means that this function needs to be a member of the
	//  FilterForm struct.

	return len(textToCheck) <= maxFilterLength
}

func (f *FilterForm) handleFilterFormInput(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	switch key {
	case tcell.KeyEnter:
		if item, _ := f.GetFocusedItemIndex(); item == filterTextField {
			f.doneHandler(key)
			return nil
		}
	}

	return event
}

func (f *FilterForm) GetText() string {
	filterText := f.filterText.GetText()
	if filterText == "" {
		return filterText
	}

	currentOptionIndex, _ := f.filterMethod.GetCurrentOption()
	globlessFilterText := strings.ReplaceAll(filterText, "*", "")

	switch currentOptionIndex {
	case filterMethodBeginsWith:
		return globlessFilterText + "*"
	case filterMethodEndsWith:
		return "*" + globlessFilterText
	case filterMethodContains:
		return "*" + globlessFilterText + "*"
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
