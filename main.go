package main

import (
	"github.com/aliics/b5/parser"
	"flag"
	"log"
)

func main() {
	repl := flag.Bool("r", false, "interactive REPL mode")
	flag.Parse()
	if !*repl {
		log.Fatalln("REPL mode is currently the only supported mode")
	}
	p := parser.NewParser(*repl)
	if err := p.Parse(); err != nil {
		log.Fatalln(err)
	}
	if err := p.Exec(); err != nil {
		log.Fatalln(err)
	}
}
