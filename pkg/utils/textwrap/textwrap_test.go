package textwrap

import "testing"

func TestShorten(t *testing.T) {
	testCases := []struct {
		text     string
		maxLen   int
		expected string
	}{
		{"Hello, world!", 10, "Hello,..."},
		{"Hello, world!", 5, "Hello..."},
		{"Hello", 10, "Hello"},
	}

	for _, tc := range testCases {
		result := Shorten(tc.text, tc.maxLen)
		if result != tc.expected {
			t.Errorf("Shorten(%q, %d) = %q; want %q", tc.text, tc.maxLen, result, tc.expected)
		}
	}
}
