package b5

type program interface {
	Exec(string) error
}

type ShellProgram struct{}

func (s ShellProgram) Exec(str string) error {
	panic("implement me")
}

type StdProgram struct{}

func (s StdProgram) Exec(str string) error {
	panic("implement me")
}
