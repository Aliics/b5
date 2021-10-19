package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/aliics/b5"
	"os"
)

func main() {
	if len(flag.Args()) == 0 {
		p := b5.ShellProgram{}

		b := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("> ")

			line, err := b.ReadString('\n')
			if err != nil {
				panic(err)
			}

			err = p.Exec(line)
			if err != nil {
				panic(err)
			}
		}
	} else {
		str, err := os.ReadFile(flag.Arg(0))
		if err != nil {
			panic(err)
		}

		p := b5.StdProgram{}
		err = p.Exec(string(str))
		if err != nil {
			return
		}
	}
}
