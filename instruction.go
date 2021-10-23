package b5

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

type ifInstruction struct {
	conditionValue  int
	executableBlock func()
}

func (i ifInstruction) exec() {
	if i.conditionValue == 1 {
		i.executableBlock()
	}
}

func createInstructions(pts []pToken) (is []instruction, err error) {
	for i := 0; i < len(pts); i++ {
		pt := pts[i]
		switch pt.tt {
		case letK:
			i, err = traverseSpaces(i, pts, true)
			if err != nil {
				return nil, err
			}

			identI := i
			if pts[i].tt != ident {
				return nil, errors.New(`expected identifier after "let"`)
			}

			i, err = traverseSpaces(i, pts, false)
			if err != nil {
				return nil, err
			}

			if pts[i].tt != equals {
				return nil, errors.New(`expected equals after identifier "` + pts[identI].data.(string) + `"`)
			}

			i, err = traverseSpaces(i, pts, false)
			if err != nil {
				return nil, err
			}

			var exp interface{}
			i, exp, err = resolveExpression(i, pts)
			if err != nil {
				return nil, err
			}

			is = append(is, letInstruction{pts[identI].data.(string), exp})
		case ifK:
			i, err = traverseSpaces(i, pts, true)
			if err != nil {
				return nil, err
			}

			var exp interface{}
			i, exp, err = resolveExpression(i+2, pts)
			if err != nil {
				return nil, err
			}

			i, err = traverseSpaces(i, pts, true)
			if err != nil {
				return nil, err
			}

			if pts[i].tt != thenK {
				return nil, errors.New(`expected "then" after if condition`)
			}

			is = append(is, ifInstruction{exp.(int), func() {
				fmt.Println("foo")
			}})
		case printF:
			i, err = traverseSpaces(i, pts, true)
			if err != nil {
				return nil, err
			}

			if pts[i].tt != ident && (pts[i].tt < stringL || pts[i].tt > numberL) {
				return nil, errors.New(`"print" requires exactly one argument`)
			}

			var exp interface{}
			i, exp, err = resolveExpression(i, pts)
			if err != nil {
				return nil, err
			}

			is = append(is, printInstruction{exp})
		}
	}

	return
}

func resolveExpression(from int, pts []pToken) (int, interface{}, error) {
	var v interface{}
	if pts[from].tt == ident {
		v = valueDefinitions[currentStack][pts[from].data.(string)]
	} else {
		v = pts[from].data
	}

	from, err := traverseSpaces(from, pts, false)
	if err != nil {
		return -1, nil, err
	}

	opI := from
	if pts[from].tt >= plus && pts[from].tt <= div {
		from, err = traverseSpaces(from, pts, false)
		if err != nil {
			return -1, nil, err
		}

		var exp interface{}
		from, exp, err = resolveExpression(from, pts)
		if err != nil {
			return -1, nil, err
		}

		expKind := reflect.TypeOf(exp).Kind()

		switch v.(type) {
		case int:
			if expKind == reflect.String {
				return -1, nil, errors.New("cannot add string to number")
			}

			switch pts[opI].tt {
			case plus:
				v = v.(int) + exp.(int)
			case minus:
				v = v.(int) - exp.(int)
			case mul:
				v = v.(int) * exp.(int)
			case div:
				v = v.(int) / exp.(int)
			}
		case string:
			switch pts[opI].tt {
			case plus:
				var str string
				if expKind == reflect.String {
					str = exp.(string)
				} else {
					str = strconv.Itoa(exp.(int))
				}

				v = v.(string) + str
			case minus:
				return -1, nil, errors.New("cannot subtract from string")
			case mul:
				if expKind != reflect.Int {
					return -1, nil, errors.New("expected number on right hand of string multiplication")
				}

				v = strings.Repeat(v.(string), exp.(int))
			case div:
				return -1, nil, errors.New("cannot divide from string")
			}
		}
	}

	return from, v, nil
}

func traverseSpaces(from int, pts []pToken, strict bool) (int, error) {
	i := from + 1
	if len(pts) <= i || len(pts) > i && pts[i].tt != space {
		if strict {
			return -1, errors.New(`expected space after "` + pts[from].tt.String() + `"`)
		} else {
			return i, nil
		}
	}

	for ; i < len(pts); i++ {
		if pts[i].tt != space {
			break
		}
	}

	return i, nil
}
