package lexer

import "github.com/shoriwe/gplasma/pkg/errors"

func (lexer *Lexer) tokenizeBinary(letterB rune) ([]rune, Kind, DirectValue, *errors.Error) {
	result := []rune{'0', letterB}
	if !lexer.reader.HasNext() {
		return result, Literal, InvalidDirectValue, errors.NewUnknownTokenKindError(lexer.line)
	}
	nextDigit := lexer.reader.Char()
	if !(nextDigit == '0' || nextDigit == '1') {
		return result, Literal, InvalidDirectValue, errors.NewUnknownTokenKindError(lexer.line)
	}
	lexer.reader.Next()
	result = append(result, nextDigit)
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if !(nextDigit == '0' || nextDigit == '1') && nextDigit != '_' {
			return result, Literal, BinaryInteger, nil
		}
		result = append(result, nextDigit)
	}
	return result, Literal, BinaryInteger, nil
}
