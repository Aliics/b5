package b5

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	currentStack     int
	valueDefinitions []map[string]interface{}
)

func interpret(pts []pToken) error {
	_, err := interpretUntil(pts, none)
	return err
}

func interpretUntil(pts []pToken, until tokenType) (i int, err error) {
	for ; i < len(pts); i++ {
		pt := pts[i]

		switch pt.tt {
		case until:
			return i, nil
		case letK:
			i, err = traverseSpaces(i, pts, true)
			if err != nil {
				return -1, err
			}

			identI := i
			if pts[i].tt != ident {
				return -1, errors.New(`expected identifier after "let"`)
			}

			i, err = traverseSpaces(i, pts, false)
			if err != nil {
				return -1, err
			}

			if pts[i].tt != equals {
				return -1, errors.New(`expected equals after identifier "` + pts[identI].data.(string) + `"`)
			}

			i, err = traverseSpaces(i, pts, false)
			if err != nil {
				return -1, err
			}

			var exp interface{}
			i, exp, err = resolveExpression(i, pts)
			if err != nil {
				return -1, err
			}

			if len(valueDefinitions) <= currentStack {
				valueDefinitions = append(valueDefinitions, make(map[string]interface{}))
			}
			valueDefinitions[currentStack][pts[identI].data.(string)] = exp
		case ifK:
			i, err = traverseSpaces(i, pts, true)
			if err != nil {
				return -1, err
			}

			var exp interface{}
			i, exp, err = resolveExpression(i, pts)
			if err != nil {
				return -1, err
			}

			i, err = traverseSpaces(i, pts, true)
			if err != nil {
				return -1, err
			}

			if pts[i].tt != thenK {
				return -1, errors.New(`expected "then" after if condition`)
			}

			i, err = traverseSpaces(i, pts, false)
			if err != nil && pts[i].tt != newline {
				return -1, errors.New(`expected space or new line after "then"`)
			}

			if exp.(int) == 1 {
				err = executeInScope(func() error {
					scopeBegin := i
					i, err = interpretUntil(pts[scopeBegin:], endK)
					i += scopeBegin
					return err
				})
				if err != nil {
					return -1, err
				}
			} else {
				required := 1
				for ; i < len(pts); i++ {
					if pts[i].tt == ifK {
						required++
					}
					if pts[i].tt == endK {
						required--
						if required == 0 {
							break
						}
					}
				}
			}
		case printF:
			i, err = traverseSpaces(i, pts, true)
			if err != nil {
				return -1, err
			}

			if pts[i].tt != ident && (pts[i].tt < stringL || pts[i].tt > numberL) {
				return -1, errors.New(`"print" requires exactly one argument`)
			}

			var exp interface{}
			i, exp, err = resolveExpression(i, pts)
			if err != nil {
				return -1, err
			}

			fmt.Println(exp)
		}
	}

	return i, nil
}

func executeInScope(f func() error) error {
	currentStack++
	valueDefinitions = append(valueDefinitions, make(map[string]interface{}))

	err := f()
	if err != nil {
		return err
	}

	valueDefinitions = valueDefinitions[0 : len(valueDefinitions)-1]
	currentStack--

	return nil
}

func resolveExpression(from int, pts []pToken) (int, interface{}, error) {
	var v interface{}
	if pts[from].tt == ident {
		v = valueDefinitions[currentStack][pts[from].data.(string)]
	} else {
		v = pts[from].data
	}

	if len(pts) <= from+2 {
		return from, v, nil
	}

	beforeSpaceI := from
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
	} else {
		return beforeSpaceI, v, err
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
