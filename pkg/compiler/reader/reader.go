package reader

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
