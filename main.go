package main

import (
	"github.com/aliics/b5/parser"
	"log"
	"os"
)

func main() {
	repl := len(os.Args) < 2
	p := parser.NewParser(repl)
	if err := p.Parse(); err != nil {
		log.Fatalln(err)
	} else if repl {
		return // repl mode calls Exec internally
	}
	if err := p.Exec(); err != nil {
		log.Fatalln(err)
	}
}
