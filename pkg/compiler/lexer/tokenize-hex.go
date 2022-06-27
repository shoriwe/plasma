package lexer

import "github.com/shoriwe/gplasma/pkg/errors"

func (lexer *Lexer) tokenizeHexadecimal(letterX rune) ([]rune, Kind, DirectValue, *errors.Error) {
	result := []rune{'0', letterX}
	if !lexer.reader.HasNext() {
		return result, Literal, InvalidDirectValue, errors.NewUnknownTokenKindError(lexer.line)
	}
	nextDigit := lexer.reader.Char()
	if !(('0' <= nextDigit && nextDigit <= '9') ||
		('a' <= nextDigit && nextDigit <= 'f') ||
		('A' <= nextDigit && nextDigit <= 'F')) {
		return result, Literal, InvalidDirectValue, errors.NewUnknownTokenKindError(lexer.line)
	}
	lexer.reader.Next()
	result = append(result, nextDigit)
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if !(('0' <= nextDigit && nextDigit <= '9') || ('a' <= nextDigit && nextDigit <= 'f') || ('A' <= nextDigit && nextDigit <= 'F')) && nextDigit != '_' {
			return result, Literal, HexadecimalInteger, nil
		}
		result = append(result, nextDigit)
	}
	return result, Literal, HexadecimalInteger, nil
}
