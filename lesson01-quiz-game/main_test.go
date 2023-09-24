package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loadQuizFile(t *testing.T) {
	var tests = []struct {
		testName                  string
		csvPath                   string
		errorExpected             bool
		expectedNumberOfQuestions int
	}{
		{testName: "invalid-csv-path", csvPath: "invalid-path.csv", errorExpected: true, expectedNumberOfQuestions: 0},
		{testName: "valid-csv", csvPath: "./testdata/valid.csv", errorExpected: false, expectedNumberOfQuestions: 3},
		{testName: "invalid-single-record-should-be-excluded", csvPath: "./testdata/bad-record.csv", errorExpected: false, expectedNumberOfQuestions: 2},
		{testName: "empty-csv", csvPath: "./testdata/empty.csv", errorExpected: false, expectedNumberOfQuestions: 0},
		{testName: "too-many-fields", csvPath: "./testdata/too-many-fields.csv", errorExpected: false, expectedNumberOfQuestions: 0},
	}

	for _, e := range tests {
		questions, err := loadQuizFile(e.csvPath)

		if e.errorExpected {
			assert.Error(t, err, fmt.Sprintf("%s: expected an error when trying to open CSV", e.testName))
		}

		if !e.errorExpected {
			assert.NoError(t, err, fmt.Sprintf("%s: expected no error when trying to open CSV", e.testName))
		}

		assert.Equal(t, e.expectedNumberOfQuestions, len(questions), fmt.Sprintf("%s: unexpected number of records returned", e.testName))
	}
}

func Test_checkAnswer(t *testing.T) {
	var tests = []struct {
		testName         string
		question         quizQuestion
		userInput        string
		questionNumber   int
		expectedResponse bool
	}{
		{testName: "correct-answer", question: quizQuestion{question: "1+1", answer: "2"}, userInput: "2\n", questionNumber: 0, expectedResponse: true},
		{testName: "ignores-whitespace", question: quizQuestion{question: "2+2", answer: " 4 "}, userInput: "4\n", questionNumber: 1, expectedResponse: true},
		{testName: "ignores-case", question: quizQuestion{question: "are you male?", answer: "yes"}, userInput: "YES\n", questionNumber: 2, expectedResponse: true},
		{testName: "ignores-case-and-whitespace", question: quizQuestion{question: "are you male?", answer: " yes "}, userInput: " YES \n", questionNumber: 3, expectedResponse: true},
		{testName: "incorrect-answer", question: quizQuestion{question: "1+1", answer: "2"}, userInput: "abc\n", questionNumber: 4, expectedResponse: false},
	}

	for _, e := range tests {
		var buf bytes.Buffer
		buf.WriteString(e.userInput)

		result := checkAnswer(e.question, e.questionNumber, &buf)

		assert.Equal(t, e.expectedResponse, result, fmt.Sprintf("%s: unexpected response from checkAnswer", e.testName))
	}
}

func Test_waitForPrompt(t *testing.T) {
	// Test user input and check what is being written by the app to stdout
	oldOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	var buf bytes.Buffer
	userInput := "\n"
	duration := "30s"

	buf.WriteString(userInput)
	err := waitForPrompt(duration, &buf)
	assert.NoError(t, err, "expected no error calling waitForPrompt")

	_ = w.Close()
	os.Stdout = oldOut
	data, _ := io.ReadAll(r)

	assert.Containsf(t, string(data), duration, "expected duration %s to be written in the console output", duration)
}

func Test_randomiseQuestions(t *testing.T) {
	original := []quizQuestion{{question: "1+1", answer: "2"}, {question: "2+2", answer: "4"}, {question: "3+3", answer: "6"}}

	randomised := randomiseQuestions(original)

	assert.Equal(t, len(original), len(randomised), "expect randomised slice to be equal length as the original")

	assert.True(t, compareSlices(original, randomised))
}

// compareSlices is a helper function to check that all the elements in original slice appear in the target slice
func compareSlices(original, target []quizQuestion) bool {
	if len(original) != len(target) {
		return false
	}

	for _, q := range original {
		found := false
		for i2, t := range target {
			if q == t {
				// found element. pop and break
				copy(target[i2:], target[i2+1:])       // Shift left one index.
				target[len(target)-1] = quizQuestion{} // Erase last element (write zero value).
				target = target[:len(target)-1]        // Truncate slice.

				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}
