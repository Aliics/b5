package parser

import (
	"strconv"
	"strings"
)

func validString(v []string) bool {
	s := strings.Join(v, " ")
	return strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`) && strings.Count(s, `"`) == 2
}

func isNum(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
