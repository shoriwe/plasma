package lexer

import "github.com/shoriwe/gplasma/pkg/errors"

func (lexer *Lexer) tokenizeComment() ([]rune, Kind, *errors.Error) {
	var content []rune
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		char := lexer.reader.Char()
		if char == '\n' {
			break
		}
		content = append(content, char)
	}
	return content, Comment, nil
}
