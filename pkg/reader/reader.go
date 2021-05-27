package reader

import (
	"bufio"
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

type FileReader struct {
	fileHandler       *os.File
	reader            io.RuneReader
	index             int
	currentChar       rune
	currentCharLength int
	finish            bool
}

func (f *FileReader) Next() {
	rune_, runeSize, readingError := f.reader.ReadRune()
	if readingError != nil {
		if readingError != io.EOF {
			panic(readingError)
		}
		f.finish = true
		return
	}
	f.currentChar = rune_
	f.currentCharLength = runeSize
	f.index += runeSize
}

func (f *FileReader) Redo() {
	if f.index > 0 {
		_, seekError := f.fileHandler.Seek(int64(f.index-f.currentCharLength), io.SeekStart)
		f.reader = bufio.NewReader(f.fileHandler)
		if seekError != nil {
			panic(seekError)
		}
		f.index -= f.currentCharLength
		f.finish = false
		f.Next()
	}
}

func (f *FileReader) HasNext() bool {
	return !f.finish
}

func (f *FileReader) Index() int {
	return f.index
}

func (f *FileReader) Char() rune {
	return f.currentChar
}

func NewFileReader(fileHandler *os.File) *FileReader {
	return &FileReader{
		fileHandler:       fileHandler,
		reader:            bufio.NewReader(fileHandler),
		index:             0,
		currentChar:       0,
		currentCharLength: 0,
		finish:            false,
	}
}
