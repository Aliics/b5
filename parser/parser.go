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
}

func NewParser(repl bool) *Parser {
	return &Parser{repl, false, 0, make([]op, 0)}
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
			default:
				if ct == "" {
					return fmt.Errorf("expected token got %s", word)
				}
				switch ct {
				case output:
					s := joinWords(arg, word)
					if validString(s) {
						p.ops = append(p.ops, op{ct, []string{s}})
						ct = ""
						arg = ""
					} else if buildingString(s) {
						arg = s
					} else {
						return fmt.Errorf("invalid string")
					}
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
	for i := p.cursor; i < len(p.ops); i++ {
		op := p.ops[i]
		switch op.t {
		case exit:
			p.stop = true
		case output:
			fmt.Println(op.args[0])
		}
		p.cursor++
	}
	return nil
}

type token string

const (
	exit   token = "STOP"
	output token = "PRINT"
)

type op struct {
	t    token
	args []string
}
