package parser

type state map[string]interface{}

func (s state) has(x string) bool {
	return s[x] != nil
}
