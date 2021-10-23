package b5

import (
	"bufio"
	"fmt"
	"os"
)

type Program interface {
	Exec() error
}

type ShellProgram struct{}

func (s ShellProgram) Exec() error {
	b := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")

		line, err := b.ReadString('\n')
		if err != nil {
			return err
		}

		pts, err := parseTokens(line)
		if err != nil {
			return err
		}

		err = interpret(pts)
		if err != nil {
			return err
		}
	}
}

type StdProgram struct{
	File string
}

func (s StdProgram) Exec() error {
	fd, err := os.ReadFile(s.File)
	if err != nil {
		return err
	}

	pts, err := parseTokens(string(fd))
	if err != nil {
		return err
	}

	return interpret(pts)
}
