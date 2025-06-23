package main

import (
	"fmt"
	"testing"
)

import "unicode"

func isUpper(text string) bool {
	for _, ch := range text {
		if !unicode.IsUpper(ch) && unicode.IsLetter(ch) {
			return false
		}
	}
	return true
}

func TestFunction(t *testing.T) {
	var tests = []struct {
		text string
		want bool
	}{
		{"UPPERCASE", true},
		{"lowercase", false},
		{"WITH SPACE", true},
		{"FOo", false},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("text: %s", test.text)
		t.Run(testName, func(t *testing.T) {
			got := isUpper(test.text)
			if got != test.want {
				t.Errorf("got %t, want %t", got, test.want)
			}
		})
	}
}
