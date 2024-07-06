package main

import "testing"

func Test_camelcase(t *testing.T) {
	tests := []struct {
		input         string
		numberOfWords int32
	}{
		{input: "one", numberOfWords: 1},
		{input: "oneTwo", numberOfWords: 2},
		{input: "oneTwoThree", numberOfWords: 3},
		{input: "x", numberOfWords: 1},
		{input: "xYZ", numberOfWords: 3},
		{input: "saveChangesInTheEditor", numberOfWords: 5},
		{input: "thisIsACamelCaseString", numberOfWords: 6},
	}

	for _, tc := range tests {
		if camelcase(tc.input) != tc.numberOfWords {
			t.Fatalf("expected %d words but got %d", tc.numberOfWords, camelcase(tc.input))
		}
	}
}
