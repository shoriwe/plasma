package lexer

import "github.com/shoriwe/gplasma/pkg/errors"

func (lexer *Lexer) tokenizeScientificFloat(base []rune) ([]rune, Kind, DirectValue, *errors.Error) {
	result := base
	if !lexer.reader.HasNext() {
		return result, Literal, InvalidDirectValue, errors.NewUnknownTokenKindError(lexer.line)
	}
	direction := lexer.reader.Char()
	if (direction != '-') && (direction != '+') {
		return result, Literal, InvalidDirectValue, errors.NewUnknownTokenKindError(lexer.line)
	}
	lexer.reader.Next()
	// Ensure next is a number
	if !lexer.reader.HasNext() {
		return result, Literal, InvalidDirectValue, errors.NewUnknownTokenKindError(lexer.line)
	}
	result = append(result, direction)
	nextDigit := lexer.reader.Char()
	if !('0' <= nextDigit && nextDigit <= '9') {
		return result, Literal, InvalidDirectValue, errors.NewUnknownTokenKindError(lexer.line)
	}
	result = append(result, nextDigit)
	lexer.reader.Next()
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if !('0' <= nextDigit && nextDigit <= '9') && nextDigit != '_' {
			return result, Literal, ScientificFloat, nil
		}
		result = append(result, nextDigit)
	}
	return result, Literal, ScientificFloat, nil
}
