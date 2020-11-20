package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Parser struct {
	repl   bool
	stop   bool
	cursor int
	ops    []op
	state  state
}

func NewParser(repl bool) *Parser {
	return &Parser{repl, false, 0, make([]op, 0), make(map[string]interface{})}
}

func (p *Parser) Parse() error {
	scanner := bufio.NewScanner(os.Stdin)
	for line := 1; !p.stop; line++ {
		fmt.Printf("> ")
		if !scanner.Scan() {
			return nil
		}
		var ct token
		var e *expr
		buildExpr := func(word string) {
			e.add(word)
			if e.isComplete() {
				p.ops = append(p.ops, op{ct, *e})
				ct = ""
				e = nil
			}
		}
		for _, word := range strings.Split(scanner.Text(), " ") {
			switch token(word) {
			case exit:
				p.ops = append(p.ops, op{token: exit})
			case output:
				ct = output
				e = newExpr(p.state)
			case variable:
				ct = variable
				e = newAssignExpr(p.state)
			default:
				if ct == "" {
					return fmt.Errorf("[%d] expected token got %s", line, word)
				}
				switch ct {
				case output, variable:
					buildExpr(word)
				}
			}
		}
		if ct != "" || e != nil {
			return fmt.Errorf("[%d] incomplete expression near %s", line, ct)
		}
		if p.repl {
			if e := p.Exec(); e != nil {
				return e
			}
		}
	}
	return nil
}

func (p *Parser) Exec() error {
	for ; p.cursor < len(p.ops); p.cursor++ {
		op := p.ops[p.cursor]
		switch op.token {
		case exit:
			p.stop = true
		case output:
			fmt.Println(op.expr.value())
		case variable:
			name := op.expr.words[0]
			if p.state.has(name) {
				return fmt.Errorf("%s is already assigned", name)
			}
			p.state[name] = op.expr.value()
			if p.repl {
				fmt.Printf("$%s = %v\n", name, p.state[name])
			}
		}
	}
	return nil
}

type token string

const (
	exit     token = "STOP"
	output   token = "PRINT"
	variable token = "LET"
)

type op struct {
	token token
	expr  expr
}
