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
	Char() uint8
}

type StringReader struct {
	content string
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

func (s *StringReader) Char() uint8 {
	return s.content[s.index]
}

func NewStringReader(code string) *StringReader {
	return &StringReader{
		content: code,
		index:   0,
		length:  len(code),
	}
}

type FileReader struct {
	fileHandler *os.File
	index       int
	currentChar uint8
	finish      bool
}

func (f *FileReader) Next() {
	charSlice := make([]byte, 1)
	_, readingError := f.fileHandler.Read(charSlice)
	if readingError != nil {
		if readingError == io.EOF {
			f.finish = true
			return
		}
	}
	f.index++
	f.currentChar = charSlice[0]
}

func (f *FileReader) Redo() {
	_, seekError := f.fileHandler.Seek(-1, io.SeekCurrent)
	if seekError != nil {
		panic(seekError)
	}
	f.index--
	f.Next()
}

func (f *FileReader) HasNext() bool {
	return !f.finish
}

func (f *FileReader) Index() int {
	return f.index
}

func (f *FileReader) Char() uint8 {
	return f.currentChar
}

func NewFileReader(fileHandler *os.File) *FileReader {
	return &FileReader{
		fileHandler: fileHandler,
		index:       -1,
		currentChar: 0,
	}
}
