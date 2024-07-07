package main

import (
	"fmt"
	"regexp"
	"strings"
)

// TODO: fix when we rotate more than the alphabet length. Currently going out of index array

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func caesarCipher(s string, k int32) string {
	isUppercase := regexp.MustCompile(`^[A-Z]$`)

	alphabetLength := int32(len(alphabet))
	rotationFactor := k

	// Original alphabet in order
	original := strings.Split(alphabet, "")

	// Stores which index in the alphabet each character has been mapped to after rotation has been applied
	mappings := make(map[string]int32)
	for idx, letter := range original {

		proposedIndex := int32(idx) + rotationFactor

		// No wrap around required
		if proposedIndex < alphabetLength {
			//fmt.Printf("Letter: %s (current index: %d, proposed index: %d)\n", letter, idx, proposedIndex)
			mappings[letter] = proposedIndex
			continue
		}

		// Index has wrapped around to the beginning
		extraPositions := proposedIndex - alphabetLength
		mappings[letter] = extraPositions
		//fmt.Printf("Letter: %s (current index: %d, proposed index: %d (wrapped))\n", letter, idx, extraPositions)
	}

	var cipherText = ""
	for _, letter := range s {
		currentLetter := string(letter)

		// Check for any letter characters and convert. Input letters can be uppercase
		if contains(original, strings.ToLower(currentLetter)) {
			// If input is uppercase return as uppercase
			if isUppercase.MatchString(currentLetter) {
				fmt.Printf("Letter %s has been swapped to %s\n", currentLetter, strings.ToUpper(original[mappings[strings.ToLower(currentLetter)]]))
				cipherText += strings.ToUpper(original[mappings[strings.ToLower(currentLetter)]])
				continue
			}
			//fmt.Printf("Letter %s has been swapped to %s\n", currentLetter, original[mappings[strings.ToLower(currentLetter)]])
			cipherText += original[mappings[currentLetter]]
			continue
		}

		// Non-alphabet character such as hyphens
		cipherText += currentLetter
	}

	return cipherText
}

func contains(s []string, str string) bool {
	var found bool
	for _, v := range s {
		if v == str {
			found = true
		}
	}
	return found
}

func main() {
	inputText := "middle-Outz"
	fmt.Printf("Input: %s\nReturn: %s\n", inputText, caesarCipher(inputText, 2))
}
