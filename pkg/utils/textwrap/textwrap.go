package textwrap

import "unicode/utf8"

func Shorten(text string, maxLen int) string {
	lastSpaceIx := maxLen
	len := 0
	for i, r := range text {
		if utf8.RuneLen(r) > 1 {
			len += 2
		} else {
			len++
		}
		if len > maxLen {
			return text[:lastSpaceIx] + "..."
		}
		if r == ' ' {
			lastSpaceIx = i
		}
	}
	return text
}
