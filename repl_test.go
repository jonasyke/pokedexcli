package main

import (
	"testing"
)


func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "   hello   world    ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "HeLLo  WoRLD!\t\n foo",
			expected: []string{"hello", "world!", "foo"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "   ",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("input %q → length mismatch: got %d, want %d",
				c.input, len(actual), len(c.expected))
		}

		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("input %q → at index %d: got %q, want %q",
					c.input, i, actual[i], c.expected[i])
			}
		}
	}
}