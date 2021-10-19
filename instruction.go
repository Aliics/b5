package b5

import (
	"errors"
	"fmt"
)

var valueDefinitions []map[string]string

type instruction interface {
	exec()
}

type printInstruction struct {
	data string
}

func (p printInstruction) exec() {
	fmt.Println(p.data)
}

type letInstruction struct {
	key, value string
}

func (p letInstruction) exec() {
	if len(valueDefinitions) < 1 {
		valueDefinitions = append(valueDefinitions, make(map[string]string))
	}
	valueDefinitions[len(valueDefinitions)-1][p.key] = p.value
}

func createInstructions(pts []pToken) (is []instruction, err error) {
	for i := 0; i < len(pts); i++ {
		pt := pts[i]
		switch pt.tt {
		case letK:
			if len(pts) < i+1 || pts[i+1].tt != space {
				return nil, errors.New(`expected space after "let"`)
			}
			if len(pts) < i+2 || pts[i+2].tt != ident {
				return nil, errors.New(`expected identifier after "let"`)
			}
			if len(pts) < i+3 || pts[i+3].tt != equals {
				return nil, errors.New(`expected equals after identifier "`+pts[i+2].data+`"`)
			}

			is = append(is, letInstruction{pts[i+2].data, pts[i+4].data})
			i += 3
		case printF:
			if len(pts) < i+1 || pts[i+1].tt != space {
				return nil, errors.New(`expected space after "print"`)
			}
			if len(pts) < i+2 || pts[i+2].tt != ident && pts[i+2].tt != stringL {
				return nil, errors.New(`"print" requires exactly one argument`)
			}

			rh := pts[i+2]

			var data string
			if rh.tt == ident {
				data = valueDefinitions[len(valueDefinitions)-1][rh.data]
			} else {
				data = rh.data
			}

			is = append(is, printInstruction{data})
			i += 1
		}
	}

	return
}
