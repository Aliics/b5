package b5

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

type pToken struct {
	tt   tokenType
	data interface{}
}

type tokenType uint8

const (
	// Special characters
	newline tokenType = iota
	space
	equals
	assert
	plus
	minus
	mul
	div
	ident
	// Keywords
	remK
	letK
	ifK
	thenK
	elseK
	endK
	// Built-in functions
	printF
	// Literals
	stringL
	numberL
)

func (t tokenType) String() string {
	switch t {
	case newline:
		return "newline"
	case space:
		return "space"
	case equals:
		return "equals"
	case assert:
		return "assert"
	case plus:
		return "plus"
	case minus:
		return "minus"
	case mul:
		return "mul"
	case div:
		return "div"
	case ident:
		return "ident"
	case remK:
		return "rem"
	case letK:
		return "let"
	case ifK:
		return "if"
	case thenK:
		return "then"
	case elseK:
		return "else"
	case endK:
		return "end"
	case printF:
		return "print"
	case stringL:
		return "string"
	case numberL:
		return "number"
	default:
		return "none"
	}
}

func parseTokens(str string) (pts []pToken, err error) {
	ln := 1
	mkErr := func(msg string) error {
		return errors.New("[" + strconv.Itoa(ln) + "]: " + msg)
	}

	for i := 0; i < len(str); i++ {
		r := unicode.ToLower(rune(str[i]))
		switch r {
		case '\n':
			ln++
			pts = append(pts, pToken{tt: newline})
		case ' ', '\t':
			pts = append(pts, pToken{tt: space})
		case '=': // =, ==
			if len(str) > i+1 && unicode.ToLower(rune(str[i+1])) == '=' {
				pts = append(pts, pToken{tt: assert})
				i++
			} else {
				pts = append(pts, pToken{tt: equals})
			}
		case '+':
			pts = append(pts, pToken{tt: plus})
		case '-':
			pts = append(pts, pToken{tt: minus})
		case '*':
			pts = append(pts, pToken{tt: mul})
		case '/':
			pts = append(pts, pToken{tt: div})
		case 'l': // LET
			if isWord(i, str, "let") {
				pts = append(pts, pToken{tt: letK})
				i += 2
			}
		case 'r': // REM
			if isWord(i, str, "rem") {
				pts = append(pts, pToken{tt: remK})
				for ; i < len(str); i++ {
					if len(str) > i+1 && str[i+1] == '\n' {
						break
					}
				}
			}
		case 'p': // PRINT
			if isWord(i, str, "print") {
				pts = append(pts, pToken{tt: printF})
				i += 4
			}
		case 'i': // IF
			if isWord(i, str, "if") {
				pts = append(pts, pToken{tt: ifK})
				i += 1
			}
		case 't': // THEN
			if isWord(i, str, "then") {
				pts = append(pts, pToken{tt: thenK})
				i += 3
			}
		case 'e': // END
			if isWord(i, str, "end") {
				pts = append(pts, pToken{tt: endK})
				i += 2
			} else if isWord(i, str, "else") {
				pts = append(pts, pToken{tt: elseK})
				i += 3
			}
		case '"': // Strings
			var end int
			for j := i + 1; j < len(str); j++ {
				if str[j] == '"' {
					end = j
					break
				}
			}

			if end == 0 {
				return nil, mkErr("string is not closed")
			}

			pts = append(pts, pToken{stringL, str[i+1 : end]})
			i = end
		default:
			if unicode.IsDigit(r) { // Numeric literals
				var end int
				for j := i; j < len(str); j++ {
					if !unicode.IsDigit(rune(str[j])) {
						end = j
						break
					}
				}

				atoi, err := strconv.Atoi(str[i:end])
				if err != nil {
					return nil, err
				}

				pts = append(pts, pToken{numberL, atoi})
				i = end - 1
			} else { // Identifiers
				var end int
				for j := i; j < len(str); j++ {
					if !unicode.IsLetter(rune(str[j])) {
						end = j
						break
					}
				}

				if end < i {
					return nil, mkErr("identifier starts with non-alphabetic rune")
				}

				pts = append(pts, pToken{ident, str[i:end]})
				i = end - 1
			}
		}
	}

	return
}

func isWord(c int, str, wanted string) bool {
	l := len(wanted)
	return len(str) >= c+l && strings.ToLower(str[c:c+l]) == wanted
}
