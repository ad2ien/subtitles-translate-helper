package internal

import (
	"testing"
)

func TestSplitLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple test",
			input:    "This is a sample translation",
			expected: "This is a sample\ntranslation",
		},
		{
			name:     "Long input test",
			input:    "This is a very long sample translation that should be split into multiple lines",
			expected: "This is a very long sample translation that\nshould be split into multiple lines",
		},
		{
			name:     "Short input test",
			input:    "Short",
			expected: "Short",
		},
		{
			name:     "Empty input test",
			input:    "",
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := splitLine(test.input)
			if actual != test.expected {
				t.Errorf("Expected %q, but got %q", test.expected, actual)
			}
		})
	}
}

func TestGetLineType(t *testing.T) {
	testCases := []struct {
		input    string
		expected LineType
	}{
		{"", SEPARATOR},
		{"123", SUB_NUMBER},
		{"00:00:10,640 --> 00:00:13,200", TIMESTAMP},
		{"♪ This is a song", SONG},
		{"This is a dialogue", DIALOGUE},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := getLineType(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}
