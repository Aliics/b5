package main

type parser struct {
	repl bool
	ts   tokens
}

func newParser(repl bool) parser {
	return parser{repl, make(tokens, 0)}
}

func (p parser) parse() error {
	return nil
}

type tokens []token

type token string

const (
	output token = "PRINT"
)
