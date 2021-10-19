package b5

type Program interface {
	Exec(string) error
}

type ShellProgram struct{}

func (s ShellProgram) Exec(str string) error {
	pts, err := parseTokens(str)
	if err != nil {
		return err
	}

	is, err := createInstructions(pts)
	if err != nil {
		return err
	}
	for _, i := range is {
		i.exec()
	}

	return nil
}

type StdProgram struct{}

func (s StdProgram) Exec(str string) error {
	panic("implement me")
}
