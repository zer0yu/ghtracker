package url

import (
	"testing"
)

func TestGetRelativeURL(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"https://github.com/AFLplusplus/LibAFL", "AFLplusplus/LibAFL"},
		{"https://github.com/AFLplusplus/LibAFL/", "AFLplusplus/LibAFL"},
		{"github.com/AFLplusplus/LibAFL", "AFLplusplus/LibAFL"},
	}

	for _, tc := range testCases {
		result, err := GetRelativeURL(tc.input)
		if err != nil {
			t.Errorf("GetRelativeURL(%q) returned error: %v", tc.input, err)
		}
		if result != tc.expected {
			t.Errorf("GetRelativeURL(%q) = %q; want %q", tc.input, result, tc.expected)
		}
	}
}
