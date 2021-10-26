package b5

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type interpreter struct {
	pts []pToken
	tkn int

	memoryManager
}

func newInterpreter(pts []pToken) *interpreter {
	return &interpreter{
		pts: pts,
		memoryManager: memoryManager{valueDefinitions: []map[string]interface{}{{}}},
	}
}

func (i *interpreter) mkErr(msg string) error {
	ln := 1
	for j := 0; j < i.tkn; j++ {
		if i.pts[j].tt == newline {
			ln++
		}
	}

	return errors.New("[" + strconv.Itoa(ln) + "]: " + msg)
}

func (i *interpreter) interpret() error {
	return i.interpretUntil([]tokenType{})
}

func (i *interpreter) interpretUntil(until []tokenType) (err error) {
	for ; i.tkn < len(i.pts); i.tkn++ {
		// Check if we go to the until tkn.
		for _, t := range until {
			if i.pts[i.tkn].tt == t {
				return nil
			}
		}

		switch i.pts[i.tkn].tt {
		case letK: // LET x = 1
			err = i.traverseSpaces(true)
			if err != nil {
				return err
			}

			identI := i.tkn
			if i.pts[i.tkn].tt != ident {
				return i.mkErr(`expected identifier after "let"`)
			}

			err = i.traverseSpaces(false)
			if err != nil {
				return err
			}

			if i.pts[i.tkn].tt != equals {
				return i.mkErr(`expected equals after identifier "` + i.pts[identI].data.(string) + `"`)
			}

			err = i.traverseSpaces(false)
			if err != nil {
				return err
			}

			exp, err := i.resolveExpression()
			if err != nil {
				return err
			}

			i.newValue(i.pts[identI].data.(string), exp)
		case ifK: // IF 1 THEN PRINT 1 ELSE PRINT 0 END
			err = i.traverseSpaces(true)
			if err != nil {
				return err
			}

			exp, err := i.resolveExpression()
			if err != nil {
				return err
			}

			err = i.traverseSpaces(true)
			if err != nil {
				return err
			}

			if i.pts[i.tkn].tt != thenK {
				return i.mkErr(`expected "then" after if condition`)
			}

			err = i.traverseSpaces(false)
			if err != nil && i.pts[i.tkn].tt != newline {
				return i.mkErr(`expected space or new line after "then"`)
			}

			if exp.(int) == 1 {
				err = i.executeInScope(func() error {
					err = i.interpretUntil([]tokenType{endK, elseK})
					return err
				})
				if err != nil {
					return err
				}

				if i.pts[i.tkn].tt == elseK {
					i.traverseToToken(endK)
				}
			} else {
				required := 1
				for ; i.tkn < len(i.pts); i.tkn++ {
					if i.pts[i.tkn].tt == ifK {
						required++
					}
					if i.pts[i.tkn].tt == endK || i.pts[i.tkn].tt == elseK {
						required--
						if required == 0 {
							break
						}
					}
				}
			}
		case printF: // PRINT 1
			err = i.traverseSpaces(true)
			if err != nil {
				return err
			}

			if i.pts[i.tkn].tt != ident && (i.pts[i.tkn].tt < stringL || i.pts[i.tkn].tt > numberL) {
				return i.mkErr(`"print" requires exactly one argument`)
			}

			exp, err := i.resolveExpression()
			if err != nil {
				return err
			}

			fmt.Println(exp)
		}
	}

	return nil
}

func (i *interpreter) resolveExpression() (interface{}, error) {
	var v interface{}
	if i.pts[i.tkn].tt == ident {
		var err error
		v, err = i.getValue(i.pts[i.tkn].data.(string))
		if err != nil {
			return nil, err
		}
	} else {
		v = i.pts[i.tkn].data
	}

	if len(i.pts) <= i.tkn+2 {
		return v, nil
	}

	beforeSpaceI := i.tkn
	err := i.traverseSpaces(false)
	if err != nil {
		return nil, err
	}

	opI := i.tkn
	if i.pts[i.tkn].tt >= assert && i.pts[i.tkn].tt <= div {
		err = i.traverseSpaces(false)
		if err != nil {
			return nil, err
		}

		exp, err := i.resolveExpression()
		if err != nil {
			return nil, err
		}

		expKind := reflect.TypeOf(exp).Kind()

		switch v.(type) {
		case int:
			if expKind == reflect.String {
				return nil, i.mkErr("cannot add string to number")
			}

			switch i.pts[opI].tt {
			case assert:
				if v.(int) == exp.(int) {
					v = 1
				} else {
					v = 0
				}
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
			switch i.pts[opI].tt {
			case assert:
				if expKind == reflect.String && v.(string) == exp.(string) {
					v = 1
				} else {
					v = 0
				}
			case plus:
				var str string
				if expKind == reflect.String {
					str = exp.(string)
				} else {
					str = strconv.Itoa(exp.(int))
				}

				v = v.(string) + str
			case minus:
				return nil, i.mkErr("cannot subtract from string")
			case mul:
				if expKind != reflect.Int {
					return nil, i.mkErr("expected number on right hand of string multiplication")
				}

				v = strings.Repeat(v.(string), exp.(int))
			case div:
				return nil, i.mkErr("cannot divide from string")
			}
		}
	} else {
		i.tkn = beforeSpaceI
		return v, err
	}

	return v, nil
}

func (i *interpreter) traverseToToken(match tokenType) {
	for ; i.tkn < len(i.pts); i.tkn++ {
		if i.pts[i.tkn].tt == match {
			break
		}
	}
}

func (i *interpreter) traverseSpaces(strict bool) error {
	i.tkn++
	if len(i.pts) <= i.tkn || len(i.pts) > i.tkn && i.pts[i.tkn].tt != space {
		if strict {
			return i.mkErr(`expected space after "` + i.pts[i.tkn].tt.String() + `"`)
		} else {
			return nil
		}
	}

	for ; i.tkn < len(i.pts); i.tkn++ {
		if i.pts[i.tkn].tt != space {
			break
		}
	}

	return nil
}

func (i interpreter) executeInScope(f func() error) error {
	i.pushStack()

	err := f()
	if err != nil {
		return err
	}

	i.popStack()

	return nil
}
