package main

import (
	"flag"
	"github.com/aliics/b5/parser"
	"log"
)

func main() {
	repl := flag.Bool("r", false, "interactive REPL mode")
	flag.Parse()
	p := parser.NewParser(*repl)
	if err := p.Parse(); err != nil {
		log.Fatalln(err)
	} else if *repl {
		return // repl mode calls Exec internally
	}
	if err := p.Exec(); err != nil {
		log.Fatalln(err)
	}
}
