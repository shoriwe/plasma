package lexer

import "github.com/shoriwe/gplasma/pkg/errors"

func (lexer *Lexer) tokenizeOctal(letterO rune) ([]rune, Kind, DirectValue, *errors.Error) {
	result := []rune{'0', letterO}
	if !lexer.reader.HasNext() {
		return result, Literal, InvalidDirectValue, errors.NewUnknownTokenKindError(lexer.line)
	}
	nextDigit := lexer.reader.Char()
	if !('0' <= nextDigit && nextDigit <= '7') {
		return result, Literal, InvalidDirectValue, errors.NewUnknownTokenKindError(lexer.line)
	}
	lexer.reader.Next()
	result = append(result, nextDigit)
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if !('0' <= nextDigit && nextDigit <= '7') && nextDigit != '_' {
			return result, Literal, OctalInteger, nil
		}
		result = append(result, nextDigit)
	}
	return result, Literal, OctalInteger, nil
}
