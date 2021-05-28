package reader

import (
	"io"
	"os"
)

type Reader interface {
	Next()
	Redo()
	HasNext() bool
	Index() int
	Char() rune
}

type StringReader struct {
	content []rune
	index   int
	length  int
}

func (s *StringReader) Next() {
	s.index++
}

func (s *StringReader) Redo() {
	s.index--
}

func (s *StringReader) HasNext() bool {
	return s.index < s.length
}

func (s *StringReader) Index() int {
	return s.index
}

func (s *StringReader) Char() rune {
	return s.content[s.index]
}

func NewStringReader(code string) *StringReader {
	return &StringReader{
		content: []rune(code),
		index:   0,
		length:  len(code),
	}
}

func NewStringReaderFromFile(file *os.File) *StringReader {
	content, readingError := io.ReadAll(file)
	if readingError != nil {
		panic(readingError)
	}
	return NewStringReader(string(content))
}
