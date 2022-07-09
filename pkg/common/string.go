package common

var directCharEscapeValue = map[rune][]rune{
	'a': {7},
	'b': {8},
	'e': {'\\', 'e'},
	'f': {12},
	'n': {10},
	'r': {13},
	't': {9},
	'?': {'\\', '?'},
}

func Repeat(s string, times int64) string {
	result := ""
	for i := int64(0); i < times; i++ {
		result += s
	}
	return result
}
