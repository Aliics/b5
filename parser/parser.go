package parser

type Parser struct {
	repl   bool
	tokens Tokens
}

func NewParser(repl bool) *Parser {
	return &Parser{repl, make(Tokens, 0)}
}

func (p *Parser) Parse() error {
	return nil
}

func (p *Parser) Exec() error {
	return nil
}

type Tokens []Token

type Token string

const (
	output Token = "PRINT"
)
