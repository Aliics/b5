package b5

import (
	"errors"
	"fmt"
)

type instruction interface {
	exec()
}

type printInstruction struct {
	data string
}

func (p printInstruction) exec() {
	fmt.Print(p.data)
}

func createInstructions(pts []pToken) (is []instruction, err error) {
	for i := 0; i < len(pts); i++ {
		pt := pts[i]
		switch pt.tt {
		case printF:
			if len(pts) > i+1 && pts[i+1].tt != space {
				return nil, errors.New(`expected space after "print"`)
			}

			is = append(is, printInstruction{pts[i+2].data})
			i += 1
		}
	}

	return
}
