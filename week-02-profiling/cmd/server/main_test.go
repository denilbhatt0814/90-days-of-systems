package main

import "testing"

func Test_isPrime(t *testing.T) {
	var testCases = []struct {
		name     string
		input    uint
		expected bool
	}{
		{name: "test 0", input: 0, expected: false},
		{name: "test 1", input: 1, expected: false},
		{name: "test 2", input: 2, expected: true},
		{name: "test 3", input: 3, expected: true},
		{name: "test 4", input: 4, expected: false},
		{name: "test 9", input: 9, expected: false},
		{name: "test 25", input: 25, expected: false},
		{name: "test 49", input: 49, expected: false},
		{name: "test 97", input: 97, expected: true},
		{name: "test 100", input: 100, expected: false},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actual := isPrime(tt.input)
			if actual != tt.expected {
				t.Errorf("isPrime(%d): expected %v, actual %v", tt.input, tt.expected, actual)
			}
		})
	}

}

func Test_countPrimeToN(t *testing.T) {
	t.Run("countPrimeToN(100)", func(t *testing.T) {
		var input uint = 100
		var expected uint = 25
		got := countPrimeToN(input)
		if got != expected {
			t.Errorf("isPrime(%d): expected %v, actual %v", input, expected, got)
		}
	})
}
