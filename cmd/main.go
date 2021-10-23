package main

import (
	"flag"
	"github.com/aliics/b5"
)

func main() {
	flag.Parse()

	var p b5.Program
	if len(flag.Args()) == 0 {
		p = b5.ShellProgram{}
	} else {
		p = b5.StdProgram{File: flag.Arg(0)}
	}

	err := p.Exec()
	if err != nil {
		panic(err)
	}
}
