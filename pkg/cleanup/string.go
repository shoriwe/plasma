package cleanup

import (
	"bytes"
)

func ReplaceEscaped(s []byte) []byte {
	// Replace escaped literals
	for _, char := range []byte("\\'\"`") {
		s = bytes.ReplaceAll(s, []byte{'\\', char}, []byte{char})
	}
	// Replace with char based
	s = bytes.ReplaceAll(s, []byte{'\\', 'a'}, []byte{7})         // a
	s = bytes.ReplaceAll(s, []byte{'\\', 'b'}, []byte{8})         // b
	s = bytes.ReplaceAll(s, []byte{'\\', 'e'}, []byte{'\\', 'e'}) // e
	s = bytes.ReplaceAll(s, []byte{'\\', 'f'}, []byte{12})        // f
	s = bytes.ReplaceAll(s, []byte{'\\', 'n'}, []byte{10})        // n
	s = bytes.ReplaceAll(s, []byte{'\\', 'r'}, []byte{13})        // r
	s = bytes.ReplaceAll(s, []byte{'\\', 't'}, []byte{9})         // t
	s = bytes.ReplaceAll(s, []byte{'\\', '?'}, []byte{'\\', '?'}) // ?
	// Replace hex with numbers
	// ToDo: Support hex escape
	// Replace unicode with numbers
	// ToDo: Support unicode escape
	return s
}
