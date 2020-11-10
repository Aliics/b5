package main

import (
	"flag"
	"log"
)

func main() {
	repl := flag.Bool("r", false, "interactive REPL mode")
	flag.Parse()
	if !*repl {
		log.Fatalln("REPL mode is currently the only supported mode")
	}
	parser := newParser(*repl)
	if err := parser.parse(); err != nil {
		log.Fatalln(err)
	}
}
