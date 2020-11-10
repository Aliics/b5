package parser

import (
	"bufio"
	"os"
	"strings"
)

type Parser struct {
	repl   bool
	stop   bool
	tokens []token
}

func NewParser(repl bool) *Parser {
	return &Parser{repl, false, make([]token, 0)}
}

func (p *Parser) Parse() error {
	scanner := bufio.NewScanner(os.Stdin)
	for !p.stop {
		if !scanner.Scan() {
			return nil
		}
		for _, word := range strings.Split(scanner.Text(), " ") {
			switch token(word) {
			case exit:
				p.tokens = append(p.tokens, exit)
			}
		}
		if p.repl {
			p.Exec()
		}
	}
	return nil
}

func (p *Parser) Exec() error {
	for _, token := range p.tokens {
		switch token {
		case exit:
			p.stop = true
		}
	}
	return nil
}

type token string

const (
	exit   token = "STOP"
	output token = "PRINT"
)
