package main

import (
	"fmt"
	"regexp"
)

// camelcase returns how many words are in the camel case formatted string. string length is always >= 1.
func camelcase(s string) int32 {
	var numOfWords int32 = 1
	isUppercase := regexp.MustCompile(`^[A-Z]$`)

	for _, letter := range s {
		if isUppercase.MatchString(string(letter)) {
			numOfWords++
		}
	}

	return numOfWords
}

func main() {
	fmt.Printf("\nNumber of words: %d\n", camelcase("helloWorld"))
}
