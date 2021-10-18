package b5

import "fmt"

type program interface {
	Exec(string) error
}

type ShellProgram struct{}

func (s ShellProgram) Exec(str string) error {
	fmt.Println(parseTokens(str))
	return nil
}

type StdProgram struct{}

func (s StdProgram) Exec(str string) error {
	panic("implement me")
}
