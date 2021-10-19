package b5

import "fmt"

func interpretLine(pts []pToken) {
	for i := 0; i < len(pts); i++ {
		pt := pts[i]
		switch pt.tt {
		case printF:
			fmt.Print(pts[i + 2].data)
			i += 2
		}
	}
}
