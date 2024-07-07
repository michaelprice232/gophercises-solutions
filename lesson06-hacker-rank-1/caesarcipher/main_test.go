package main

import "testing"

func Test_camelcase(t *testing.T) {
	tests := []struct {
		plainText      string
		rotationNumber int32
		expected       string
	}{
		{plainText: "middle-Outz", rotationNumber: 2, expected: "okffng-Qwvb"},
		{plainText: "www.abc.xy", rotationNumber: 87, expected: "fff.jkl.gh"},
		{plainText: "159357lcfd", rotationNumber: 98, expected: "159357fwzx"},
		{plainText: "hello", rotationNumber: 1, expected: "ifmmp"},
		{plainText: "golang", rotationNumber: 26, expected: "golang"},
		{plainText: "testing", rotationNumber: 0, expected: "testing"},
		{plainText: "ABCDEFGHIJKLMNOPQRSTUVWXYZ", rotationNumber: 15, expected: "PQRSTUVWXYZABCDEFGHIJKLMNO"},
	}

	for _, tc := range tests {
		result := caesarCipher(tc.plainText, tc.rotationNumber)
		if result != tc.expected {
			t.Fatalf("expected ciphertext %s back but got %s", tc.expected, result)
		}
	}
}
