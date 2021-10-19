package b5

import (
	"errors"
	"strings"
	"unicode"
)

type pToken struct {
	tt   tokenType
	data string
}

type tokenType uint8

const (
	newline tokenType = iota
	space
	equals
	ident
	remK
	letK
	dataK
	readK
	restoreK
	printF
	stringL
)

func parseTokens(str string) (pts []pToken, err error) {
	for i := 0; i < len(str); i++ {
		r := unicode.ToLower(rune(str[i]))
		switch r {
		case '\n':
			pts = append(pts, pToken{tt: newline})
		case ' ':
			pts = append(pts, pToken{tt: space})
		case '=':
			pts = append(pts, pToken{tt: equals})
		case 'l': // LET
			if isWord(i, str, "let") {
				pts = append(pts, pToken{tt: letK})
				i += 2
			}
		case 'd': // DATA
			if isWord(i, str, "data") {
				pts = append(pts, pToken{tt: dataK})
				i += 3
			}
		case 'r': // REM, READ, RESTORE
			if isWord(i, str, "rem") {
				pts = append(pts, pToken{tt: remK})
				for ; i < len(str); i++ {
					if len(str) > i+1 && str[i+1] == '\n' {
						break
					}
				}
			} else if isWord(i, str, "read") {
				pts = append(pts, pToken{tt: readK})
				i += 3
			} else if isWord(i, str, "restore") {
				pts = append(pts, pToken{tt: restoreK})
				i += 6
			}
		case 'p': // PRINT
			if isWord(i, str, "print") {
				pts = append(pts, pToken{tt: printF})
				i += 4
			}
		case '"': // Strings
			var end int
			for j := i + 1; j < len(str); j++ {
				if str[j] == '"' {
					end = j
				}
			}

			if end == 0 {
				return nil, errors.New("string is not closed")
			}

			pts = append(pts, pToken{stringL, str[i+1 : end]})
			i = end
		default:
			var end int
			for j := i; j < len(str); j++ {
				if !unicode.IsLetter(rune(str[j])) {
					end = j
					break
				}
			}

			if end < i {
				return nil, errors.New("identifier starts with non-alphabetic rune")
			}

			pts = append(pts, pToken{ident, str[i:end]})
			i = end-1
		}
	}

	return
}

func isWord(c int, str, wanted string) bool {
	l := len(wanted)
	return len(str) >= c+l && strings.ToLower(str[c:c+l]) == wanted
}
