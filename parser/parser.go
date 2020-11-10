package parser

import (
	"bufio"
	"os"
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
		line := scanner.Text()
		switch token(line) {
		case exit:
			p.stop = true
		}
	}
	return nil
}

func (p *Parser) Exec() error {
	return nil
}

type token string

const (
	exit   token = "STOP"
	output token = "PRINT"
)
