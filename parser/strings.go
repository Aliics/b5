package parser

import (
	"strings"
)

func validString(s string) bool {
	return s[0] == '"' && s[len(s)-1:] == `"` && strings.Count(s, `"`) == 2
}

func buildingString(s string) bool {
	return s[0] == '"' && strings.Count(s, `"`) == 1
}

func joinWords(word0, word1 string) string {
	if word0 == "" {
		return word1
	} else {
		return word0 + " " + word1
	}
}
