package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"testing"
)

func Test_FilterForm_handleFilterAcceptance_DoesNotAcceptSlashes(t *testing.T) {
	filterForm := CreateFilterForm()
	lastChars := []rune{'/', '\\'}

	for _, lastChar := range lastChars {
		if filterForm.handleFilterAcceptance("", lastChar) {
			t.Errorf("Expected last char '%c' not to be accepted", lastChar)
		}
	}
}

func Test_FilterForm_handleFilterAcceptance_FilterTextDoesNotExceedMaximumLength(t *testing.T) {
	filterForm := CreateFilterForm()
	validText := "this-text-is-exactly32characters"

	if !filterForm.handleFilterAcceptance(validText, ' ') {
		t.Errorf(
			"Expected text '%v' of length %v to be accepted",
			validText,
			32)
	}

	invalidText := "this-text-is-exactly-33characters"

	if filterForm.handleFilterAcceptance(invalidText, ' ') {
		t.Errorf(
			"Expected text '%v' of length %v not to be accepted",
			invalidText,
			33)
	}
}

func Test_FilterForm_handleFilterAcceptance_DoesNotAcceptGlobCharactersUnlessInManualGlobMode(t *testing.T) {
	filterForm := CreateFilterForm()
	globChars := []rune{'*', '?', '[', ']', '!'}

	runNonGlobModeTest := func(filterMethod int) {
		filterForm.filterMethod.SetCurrentOption(filterMethod)
		for _, globChar := range globChars {
			if filterForm.handleFilterAcceptance("", globChar) {
				t.Errorf("Expected last char '%c' not to be accepted", globChar)
			}
		}
	}
	runNonGlobModeTest(filterMethodBeginsWith)
	runNonGlobModeTest(filterMethodEndsWith)
	runNonGlobModeTest(filterMethodContains)

	filterForm.filterMethod.SetCurrentOption(filterMethodGlobPattern)
	for _, globChar := range globChars {
		if !filterForm.handleFilterAcceptance("", globChar) {
			t.Errorf("Expected last char '%c' to be accepted", globChar)
		}
	}
}

func Test_FilterForm_CreateFilterForm_SetsFocusToInputField(t *testing.T) {
	screen := tcell.NewSimulationScreen("") // "" = UTF-8 charset
	app := tview.NewApplication().SetScreen(screen)
	filterForm := CreateFilterForm()

	app.SetRoot(filterForm, false)

	expectedFocus := filterForm.filterText
	result := app.GetFocus()

	if result != expectedFocus {
		t.Errorf("Expected focused component to be '%v', got '%v' instead", expectedFocus, result)
	}
}

func Test_FilterForm_Clear_SetsFilterTextInputFieldToEmptyString(t *testing.T) {
	filterForm := CreateFilterForm()
	filterForm.filterText.SetText("this should not appear after test")

	filterForm.Clear()

	result := filterForm.filterText.GetText()

	if len(result) > 0 {
		t.Errorf("Expected the filter text to be empty, got '%s' instead", result)
	}
}

func Test_FilterForm_GetText_ReturnsBeginsWithGlobPattern(t *testing.T) {
	filterForm := CreateFilterForm()
	filterForm.filterText.SetText("pattern")
	filterForm.filterMethod.SetCurrentOption(filterMethodBeginsWith)

	result := filterForm.GetText()
	expected := "pattern*"

	if result != expected {
		t.Errorf("Expected pattern to be '%s', got '%s' instead", expected, result)
	}
}

func Test_FilterForm_GetText_ReturnsEndsWithGlobPattern(t *testing.T) {
	filterForm := CreateFilterForm()
	filterForm.filterText.SetText("pattern")
	filterForm.filterMethod.SetCurrentOption(filterMethodEndsWith)

	result := filterForm.GetText()
	expected := "*pattern"

	if result != expected {
		t.Errorf("Expected pattern to be '%s', got '%s' instead", expected, result)
	}
}

func Test_FilterForm_GetText_ReturnsContainsGlobPattern(t *testing.T) {
	filterForm := CreateFilterForm()
	filterForm.filterText.SetText("pattern")
	filterForm.filterMethod.SetCurrentOption(filterMethodContains)

	result := filterForm.GetText()
	expected := "*pattern*"

	if result != expected {
		t.Errorf("Expected pattern to be '%s', got '%s' instead", expected, result)
	}
}

func Test_FilterForm_GetText_ReturnsManualGlobPattern(t *testing.T) {
	filterForm := CreateFilterForm()
	filterForm.filterText.SetText("pat*ter*n")
	filterForm.filterMethod.SetCurrentOption(filterMethodGlobPattern)

	result := filterForm.GetText()
	expected := "pat*ter*n"

	if result != expected {
		t.Errorf("Expected pattern to be '%s', got '%s' instead", expected, result)
	}
}

func Test_FilterForm_GetText_ReturnsEmptyStringWhenFilterTextFieldIsEmpty(t *testing.T) {
	filterForm := CreateFilterForm()
	filterForm.filterText.SetText("")

	filterMethods := []int{filterMethodBeginsWith, filterMethodEndsWith, filterMethodContains, filterMethodGlobPattern}

	for i := range filterMethods {
		filterForm.filterMethod.SetCurrentOption(i)

		result := filterForm.GetText()
		expected := ""

		if result != expected {
			t.Errorf("Expected pattern to be '%s', got '%s' instead", expected, result)
		}
	}
}

func Test_FilterForm_SetDoneHandler_SetsDoneFunctionOfFilterTextField(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	app := tview.NewApplication().SetScreen(screen)
	filterForm := CreateFilterForm()
	doneFuncCalled := false
	doneFunc := func(key tcell.Key) {
		doneFuncCalled = true
	}
	filterForm.SetDoneHandler(doneFunc)

	app.SetRoot(filterForm, false)

	inputHandler := filterForm.InputHandler()
	inputHandler(
		tcell.NewEventKey(tcell.KeyEnter, rune(tcell.KeyEnter), tcell.ModNone),
		func(p tview.Primitive) {})

	if !doneFuncCalled {
		t.Error("Expected done function to be set, but it was not")
	}
}

func Test_FilterForm_SetDoneHandler_ReturnsFilterFormFromFunc(t *testing.T) {
	filterForm := CreateFilterForm()

	result := filterForm.SetDoneHandler(func(key tcell.Key) {})

	if result != filterForm {
		t.Errorf("Expected the returned filter form to be '%v', got '%v' instead", filterForm, result)
	}
}

func Test_FilterForm_SetText_SetsTextOfFilterTextField(t *testing.T) {
	filterForm := CreateFilterForm()
	expected := "bananas"
	filterForm.SetText(expected)

	result := filterForm.filterText.GetText()

	if result != expected {
		t.Errorf("Expected filter text to be set to '%s', got '%s' instead", expected, result)
	}
}

func Test_FilterForm_SetText_ReturnsFilterFormFromFunc(t *testing.T) {
	filterForm := CreateFilterForm()

	result := filterForm.SetText("test")

	if result != filterForm {
		t.Errorf("Expected the returned filter form to be '%v', got '%v' instead", filterForm, result)
	}
}

func Test_FilterForm_handleFilterFormInput_RunsDoneHandlerWhenEnterIsPressedAndFilterTextFieldIsSelected(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	app := tview.NewApplication().SetScreen(screen)
	filterForm := CreateFilterForm()
	doneFuncCalled := false
	doneFunc := func(key tcell.Key) {
		doneFuncCalled = true
	}
	filterForm.SetDoneHandler(doneFunc)

	app.SetRoot(filterForm, false)

	filterForm.handleFilterFormInput(tcell.NewEventKey(tcell.KeyEnter, rune(tcell.KeyEnter), tcell.ModNone))

	if !doneFuncCalled {
		t.Error("Expected done function to be set, but it was not")
	}
}

func Test_FilterForm_handleFilterFormInput_RunsDoneHandlerWhenEscIsPressedAndFilterTextFieldIsSelected(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	app := tview.NewApplication().SetScreen(screen)
	filterForm := CreateFilterForm()
	doneFuncCalled := false
	doneFunc := func(key tcell.Key) {
		doneFuncCalled = true
	}
	filterForm.SetDoneHandler(doneFunc)

	app.SetRoot(filterForm, false)

	filterForm.handleFilterFormInput(tcell.NewEventKey(tcell.KeyEsc, rune(tcell.KeyEsc), tcell.ModNone))

	if !doneFuncCalled {
		t.Error("Expected done function to be set, but it was not")
	}
}
