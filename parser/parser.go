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
	state  map[string]interface{}
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
			if !e.isValid() {
				return
			}
			if e.isComplete() {
				p.ops = append(p.ops, op{ct, *e})
				ct = ""
				e = nil
			}
		}
		for _, word := range strings.Split(scanner.Text(), " ") {
			switch token(word) {
			case exit:
				p.ops = append(p.ops, op{t: exit})
			case output:
				ct = output
			case variable:
				ct = variable
			default:
				if ct == "" {
					return fmt.Errorf("[%d] expected token got %s", line, word)
				}
				switch ct {
				case output:
					if e == nil {
						e = newExpr(word)
						continue
					}
					buildExpr(word)
				case variable:
					if !strings.Contains(word, "=") {
						e = newExpr(word)
					} else if e != nil {
						buildExpr(word)
					}
				}
			}
		}
		if ct != "" || e != nil {
			return fmt.Errorf("[%d] incomplete expression near %s", line, ct)
		}
		if p.repl {
			p.Exec()
		}
	}
	return nil
}

func (p *Parser) Exec() error {
	for ; p.cursor < len(p.ops); p.cursor++ {
		op := p.ops[p.cursor]
		switch op.t {
		case exit:
			p.stop = true
		case output:
			fmt.Println(op.e.value())
		case variable:
			name := op.e.words[0]
			p.state[name] = op.e.words[1]
			if p.repl {
				fmt.Printf("$%v = %v\n", name, p.state[name])
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
	t    token
	e    expr
}
