package main

import (
	"regexp"
	"strings"
)

func Format(str string, size int) string {
	replaced := strings.ReplaceAll(str, "\n", "\\n")

	spaceRegex := regexp.MustCompile(`\s+`)
	replaced = spaceRegex.ReplaceAllString(replaced, " ")
	runes := []rune(replaced)
	if len(runes) > size {
		return string(runes[:size])
	}
	return string(runes)
}
