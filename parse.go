package b5

type ptoken uint8

const (
	newLine ptoken = iota
	let
	data
	read
	restore
)

func parseTokens(str string) (pt []ptoken) {
	for i := 0; i < len(str); i++ {
		r := rune(str[i])
		switch r {
		case '\n':
			pt = append(pt, newLine)
		case 'l': // LET
			if isWord(i, str, "let") {
				pt = append(pt, let)
				i += 3
			}
		case 'd': // DATA
			if isWord(i, str, "data") {
				pt = append(pt, data)
				i += 4
			}
		case 'r': // READ, RESTORE
			if isWord(i, str, "read") {
				pt = append(pt, read)
				i += 4
			} else if isWord(i, str, "restore") {
				pt = append(pt, restore)
				i += 4
			}
		}
	}

	return
}

func isWord(c int, str, wanted string) bool {
	l := len(wanted)
	return len(str) >= c+l && str[c:c+l] == wanted
}
