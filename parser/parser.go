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
	for !p.stop {
		if !scanner.Scan() {
			return nil
		}
		var ct token
		var arg string
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
					return fmt.Errorf("expected token got %s", word)
				}
				switch ct {
				case output:
					if v := p.state[word]; v != nil {
						p.ops = append(p.ops, op{ct, []string{fmt.Sprintf("%v", v)}})
						ct = ""
						arg = ""
						break
					}
					s := joinWords(arg, word)
					if validString(s) {
						p.ops = append(p.ops, op{ct, []string{s}})
						ct = ""
						arg = ""
					} else if buildingString(s) {
						arg = s
					} else {
						return fmt.Errorf("expected expression got %v", word)
					}
				case variable:
					if !strings.Contains(word, "=") {
						return fmt.Errorf("expected assignment got %v", word)
					}
					parts := strings.Split(word, "=")
					p.ops = append(p.ops, op{ct, []string{parts[0], parts[1]}})
					ct = ""
					arg = ""
				}
			}
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
			fmt.Println(op.args[0])
		case variable:
			name := op.args[0]
			p.state[name] = op.args[1]
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
	args []string
}
