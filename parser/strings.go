package parser

import (
	"strings"
)

func validString(v []string) bool {
	s := strings.Join(v, " ")
	return strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`) && strings.Count(s, `"`) == 2
}
