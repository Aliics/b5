package b5

type program interface {
	Exec(string) error
}

type ShellProgram struct{}

func (s ShellProgram) Exec(str string) error {
	pts, err := parseTokens(str)
	if err != nil {
		return err
	}

	err = validSyntax(pts)
	if err != nil {
		return err
	}

	interpretLine(pts)

	return nil
}

type StdProgram struct{}

func (s StdProgram) Exec(str string) error {
	panic("implement me")
}
