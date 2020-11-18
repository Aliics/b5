package parser

import (
	"strings"
)

type expr struct {
	words  []string
	assign bool
}

func newExpr(words ...string) *expr {
	return &expr{words, false}
}

func newAssignExpr(words ...string) *expr {
	return &expr{words, true}
}

func (e *expr) add(word string) {
	e.words = append(e.words, word)
}

func (e *expr) isComplete() bool {
	return validString(e.words) || (e.assign && len(e.words) > 2 && e.words[1] == "=")
}

func (e *expr) value() string {
	if e.assign {
		return strings.Join(e.words[2:], " ")
	}
	return strings.Join(e.words, " ")
}
