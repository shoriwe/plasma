package lexer

import "github.com/shoriwe/gplasma/pkg/errors"

func (lexer *Lexer) tokenizeWord(startingChar rune) ([]rune, Kind, DirectValue, *errors.Error) {
	content := []rune{startingChar}
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		char := lexer.reader.Char()
		if ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || ('0' <= char && char <= '9') || (char == '_') {
			content = append(content, char)
		} else {
			break
		}
	}
	kind, directValue := detectKindAndDirectValue(content)
	return content, kind, directValue, nil
}
