package b5

import (
	"errors"
	"fmt"
	"strconv"
)

var valueDefinitions []map[string]interface{}
var currentStack int // everything is in global scope for now

type instruction interface {
	exec()
}

type printInstruction struct {
	data interface{}
}

func (p printInstruction) exec() {
	fmt.Println(p.data)
}

type letInstruction struct {
	key   string
	value interface{}
}

func (p letInstruction) exec() {
	if len(valueDefinitions) <= currentStack {
		valueDefinitions = append(valueDefinitions, make(map[string]interface{}))
	}
	valueDefinitions[currentStack][p.key] = p.value
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
				return nil, errors.New(`expected equals after identifier "` + pts[i+2].data.(string) + `"`)
			}

			exp, err := resolveExpression(pts, i+4)
			if err != nil {
				return nil, err
			}

			is = append(is, letInstruction{pts[i+2].data.(string), exp})
			i += 3
		case printF:
			if len(pts) < i+1 || pts[i+1].tt != space {
				return nil, errors.New(`expected space after "print"`)
			}
			if len(pts) < i+2 || pts[i+2].tt != ident && pts[i+2].tt != stringL && pts[i+2].tt != numberL {
				return nil, errors.New(`"print" requires exactly one argument`)
			}

			exp, err := resolveExpression(pts, i+2)
			if err != nil {
				return nil, err
			}

			is = append(is, printInstruction{exp})
			i += 1
		}
	}

	return
}

func resolveExpression(pts []pToken, from int) (v interface{}, err error) {
	if pts[from].tt == ident {
		v = valueDefinitions[currentStack][pts[from].data.(string)]
	} else {
		v = pts[from].data
	}

	if len(pts) > from+2 && pts[from+1].tt == plus {
		exp, err := resolveExpression(pts, from+2)
		if err != nil {
			return nil, err
		}

		switch v.(type) {
		case int:
			if pts[from+2].tt == stringL {
				return nil, errors.New("cannot add string to number")
			}

			v = v.(int) + exp.(int)
		case string:
			var str string
			if pts[from+2].tt == stringL {
				str = exp.(string)
			} else {
				str = strconv.Itoa(exp.(int))
			}

			v = v.(string) + str
		}
	}

	return
}
