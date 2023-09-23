package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type quizQuestion struct {
	question string
	answer   string
}

func main() {
	var csvPath = flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer' (default problems.csv)")
	var limit = flag.String("limit", "30s", "the time limit for the quiz' (default 30s)")
	var random = flag.Bool("random", false, "whether to randomise the questions (default false)")
	flag.Parse()

	timeLimit, err := time.ParseDuration(*limit)
	if err != nil {
		log.Fatalf("Problem parsing the limit flag as a duration: %v", err)
	}

	correctAnswers := 0
	resultsChannel := make(chan quizQuestion)
	completedChannel := make(chan bool)

	// Parse CSV for quiz questions
	questions := loadQuizFile(*csvPath)

	if *random {
		fmt.Printf("Randomising the questions...\n")
		questions = randomiseQuestions(questions)
	}

	// Run Quiz
	waitForPrompt(*limit)
	go askQuestions(resultsChannel, completedChannel, questions)
	timeout := time.Tick(timeLimit)

	for {
		select {
		case <-timeout:
			fmt.Print("\nTimeout!")
			printResults(correctAnswers, len(questions))
			os.Exit(0)

		case <-completedChannel:
			printResults(correctAnswers, len(questions))
			os.Exit(0)

		case _, ok := <-resultsChannel:
			// ok indicates that it received an event rather than a zero value caused by the channel closing
			if ok {
				correctAnswers++
			}
		}
	}
}

// loadQuizFile loads a CSV file from csvPath and returns them as a slice a quiz questions.
func loadQuizFile(csvPath string) []quizQuestion {
	questions := make([]quizQuestion, 0)

	b, err := os.ReadFile(csvPath)
	if err != nil {
		log.Fatalf("Problem reading the quiz file at '%s': %v", csvPath, err)
	}

	reader := csv.NewReader(bytes.NewReader(b))

	for {
		csvRecord, err := reader.Read()
		if err == io.EOF {
			break
		}
		// CSV packages expects all the records to have the same number of fields. Skip any that are irregular
		if err != nil {
			continue
		}
		// Expect exactly one question and answer in the CSV record
		if len(csvRecord) != 2 {
			continue
		}

		q := quizQuestion{question: csvRecord[0], answer: csvRecord[1]}
		questions = append(questions, q)
	}
	return questions
}

// printResults prints out the number of correctly answered questions vs total quiz questions.
func printResults(correctAnswers, totalNumber int) {
	fmt.Printf("\nYou answered %d out of %d correct!", correctAnswers, totalNumber)
}

// askQuestions iterates through the questions and calls checkAnswer for each one.
// Correctly answered questions are sent to the rc channel (to be totalled by the main go routine).
// When all questions have been processed an event is sent to the cc channel to indicate completion.
func askQuestions(rc chan quizQuestion, cc chan bool, questions []quizQuestion) {
	for i, question := range questions {
		result := checkAnswer(question, i)
		// Send correct answers via the channel to be totalled in the main go routine
		if result {
			rc <- question
		}
	}
	close(rc)

	// Indicate we have processed all the quiz questions
	cc <- true
	close(cc)
}

// checkAnswer asks the user a question on the terminal and inspects the response via stdin.
// All whitespace and case or ignored when comparing answers.
func checkAnswer(question quizQuestion, number int) bool {
	fmt.Printf("Question #%d: %s = ", number+1, question.question)

	readStdin := bufio.NewReader(os.Stdin)
	// Using ReadLine does not include the trailing \n like with ReadString
	answer, _, err := readStdin.ReadLine()
	if err != nil {
		return false
	}
	// Trim all whitespace and ignore case
	if strings.TrimSpace(strings.ToLower(string(answer))) == strings.TrimSpace(strings.ToLower(question.answer)) {
		return true
	}

	return false
}

// waitForPrompt prompts the user to press any key before the quick (and timer) starts.
func waitForPrompt(duration string) {
	fmt.Printf("Enter any key to start timer (%s): ", duration)

	readStdin := bufio.NewReader(os.Stdin)
	_, _, err := readStdin.ReadLine()
	if err != nil {
		log.Fatalf("Problem reading from stdin: %v", err)
	}

	return
}

// randomiseQuestions is called when the --random flag is set, which randomises the questions from the quiz file.
func randomiseQuestions(originalSlice []quizQuestion) []quizQuestion {
	size := len(originalSlice)

	newSlice := make([]quizQuestion, size)
	newChecker := make([]bool, size)

	for i := 0; i < size; i++ {
		// find a new random index position
		for {
			targetIndex := int(rand.Int31n(int32(size)))

			if newChecker[targetIndex] {
				continue
			} else {
				newChecker[targetIndex] = true
				newSlice[targetIndex] = originalSlice[i]
				break
			}
		}
	}

	return newSlice
}
