package main

import (
	"flag"
	"fmt"
	"github.com/aliics/b5"
	"os"
)

func main() {
	if len(flag.Args()) == 0 {
		p := b5.ShellProgram{}

		var line string
		_, err := fmt.Scanln(&line)
		if err != nil {
			return
		}

		for {
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
