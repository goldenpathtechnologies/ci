package ui

import "testing"

func TestHandleFilterAcceptanceLastChar(t *testing.T) {
	lastChars := []rune{ '/', '\\' }

	for _, lastChar := range lastChars {
		if handleFilterAcceptance("", lastChar) {
			t.Errorf("Expected last char '%c' not to be accepted", lastChar)
		}
	}
}

func TestHandleFilterAcceptanceTextLength(t *testing.T) {
	validText := "this-text-is-exactly32characters"

	if !handleFilterAcceptance(validText, ' ') {
		t.Errorf(
			"Expected text '%v' of length %v to be accepted",
			validText,
			32)
	}

	invalidText := "this-text-is-exactly-33characters"

	if handleFilterAcceptance(invalidText, ' ') {
		t.Errorf(
			"Expected text '%v' of length %v not to be accepted",
			invalidText,
			33)
	}
}