package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loadQuizFile(t *testing.T) {
	var tests = []struct {
		test                      string
		csvPath                   string
		errorExpected             bool
		expectedNumberOfQuestions int
	}{
		{test: "invalid-csv-path", csvPath: "invalid-path.csv", errorExpected: true, expectedNumberOfQuestions: 0},
		{test: "valid-csv", csvPath: "./testdata/valid.csv", errorExpected: false, expectedNumberOfQuestions: 3},
		{test: "invalid-single-record-should-be-excluded", csvPath: "./testdata/bad-record.csv", errorExpected: false, expectedNumberOfQuestions: 2},
		{test: "empty-csv", csvPath: "./testdata/empty.csv", errorExpected: false, expectedNumberOfQuestions: 0},
		{test: "too-many-fields", csvPath: "./testdata/too-many-fields.csv", errorExpected: false, expectedNumberOfQuestions: 0},
	}

	for _, e := range tests {
		questions, err := loadQuizFile(e.csvPath)

		if e.errorExpected {
			assert.Error(t, err, fmt.Sprintf("%s: expected an error when trying to open CSV", e.test))
		}

		if !e.errorExpected {
			assert.NoError(t, err, fmt.Sprintf("%s: expected no error when trying to open CSV", e.test))
		}

		assert.Equal(t, e.expectedNumberOfQuestions, len(questions), fmt.Sprintf("%s: unexpected number of records returned", e.test))
	}
}
