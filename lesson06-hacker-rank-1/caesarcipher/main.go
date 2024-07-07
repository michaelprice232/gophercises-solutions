package main

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// caesarCipher returns cipher text based on the Caesar Cipher
// https://en.wikipedia.org/wiki/Caesar_cipher
func caesarCipher(s string, k int32) string {
	// Original alphabet in order
	original := strings.Split(alphabet, "")

	mappings := rotateCharacters(original, k)

	return generateCipherText(s, original, mappings)
}

// rotateCharacters rotates characters in original to the right rotationNumber of times and returns a new index mapping for each character.
func rotateCharacters(original []string, rotationNumber int32) map[string]int32 {
	alphabetLength := int32(len(alphabet))
	mappings := make(map[string]int32)

	for idx, letter := range original {
		proposedIndex := int32(idx) + rotationNumber

		// No wrap around required
		if proposedIndex < alphabetLength {
			//fmt.Printf("Letter: %s (current index: %d, proposed index: %d)\n", letter, idx, proposedIndex)
			mappings[letter] = proposedIndex
			continue
		}

		// Index has wrapped around to the beginning. It might wrap multiple times, so we use the remainder operator
		extraPositions := proposedIndex % alphabetLength
		mappings[letter] = extraPositions
		//fmt.Printf("Letter: %s (current index: %d, proposed index: %d (wrapped))\n", letter, idx, extraPositions)
	}

	return mappings
}

func generateCipherText(inputStr string, original []string, mappings map[string]int32) string {
	isUppercase := regexp.MustCompile(`^[A-Z]$`)
	var cipherText string

	for _, letter := range inputStr {
		currentLetter := string(letter)

		// Check for any letter characters and convert. Input letters can be uppercase
		if slices.Contains(original, strings.ToLower(currentLetter)) {
			// If input is uppercase return as uppercase
			if isUppercase.MatchString(currentLetter) {
				//fmt.Printf("Letter %s has been swapped to %s\n", currentLetter, strings.ToUpper(original[mappings[strings.ToLower(currentLetter)]]))
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

func main() {
	inputText := "159357lcfd"
	fmt.Printf("Input: %s\nReturn: %s\n", inputText, caesarCipher(inputText, 98))
}
