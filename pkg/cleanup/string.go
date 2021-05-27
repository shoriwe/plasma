package cleanup

import (
	"bytes"
	"regexp"
	"strconv"
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
	hexEscapes := regexp.MustCompile("(?i)\\\\[Xx][0-9a-fA-F][0-9a-fA-F]").FindAll(s, -1)
	for _, hexEscape := range hexEscapes {
		number, parsingError := strconv.ParseUint(string(hexEscape[2:]), 16, 8)
		if parsingError != nil {
			panic(parsingError)
		}
		s = bytes.ReplaceAll(s, hexEscape, []byte{uint8(number)})
	}
	// Replace unicode with numbers
	unicodeEscapes := regexp.MustCompile("(?i)\\\\[Uu][0-9a-fA-F]{4,6}").FindAll(s, -1)
	for _, unicodeEscape := range unicodeEscapes {
		number, parsingError := strconv.ParseUint(string(unicodeEscape[2:]), 16, 32)
		if parsingError != nil {
			panic(parsingError)
		}
		s = bytes.ReplaceAll(s, unicodeEscape, []byte(string(rune(number))))
	}
	return s
}
