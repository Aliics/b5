package parser

import (
	"strings"
)

type expr struct {
	words []string
}

func newExpr(words ...string) *expr {
	return &expr{words}
}

func (e *expr) add(word string) {
	e.words = append(e.words, word)
}

func (e *expr) isValid() bool {
	return strings.HasPrefix(e.words[0], `"`)
}

func (e *expr) isComplete() bool {
	return strings.HasPrefix(e.words[0], `"`) && strings.HasSuffix(e.words[len(e.words) - 1], `"`)
}

func (e *expr) value() string {
	return strings.Join(e.words, " ")
}
