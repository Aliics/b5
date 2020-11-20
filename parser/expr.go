package parser

import (
	"strings"
)

type expr struct {
	words  []string
	assign bool
	state  state
}

func newExpr(state state, words ...string) *expr {
	return &expr{words, false, state}
}

func newAssignExpr(state state, words ...string) *expr {
	return &expr{words, true, state}
}

func (e *expr) add(word string) {
	e.words = append(e.words, word)
}

func (e *expr) isComplete() bool {
	return validString(e.words) ||
		isNum(e.words[0]) ||
		(len(e.words) == 1 && e.state.has(e.words[0])) ||
		func() bool {
			assign := e.assign && len(e.words) > 2 && e.words[1] == "="
			compString := strings.HasSuffix(e.words[len(e.words)-1], `"`)
			num := len(e.words) == 3 && isNum(e.words[2])
			access := len(e.words) == 3 && e.state.has(e.words[2])
			return assign && (compString || num || access)
		}()
}

func (e *expr) value() interface{} {
	if last := e.words[len(e.words)-1]; e.state.has(last) {
		return e.state[last]
	}
	if e.assign {
		return strings.Join(e.words[2:], " ")
	}
	return strings.Join(e.words, " ")
}
