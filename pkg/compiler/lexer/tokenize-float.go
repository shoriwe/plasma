package lexer

import "github.com/shoriwe/gplasma/pkg/errors"

func (lexer *Lexer) tokenizeFloat(base []rune) ([]rune, Kind, DirectValue, *errors.Error) {
	if !lexer.reader.HasNext() {
		lexer.reader.Redo()
		return base[:len(base)-1], Literal, Integer, nil
	}
	nextDigit := lexer.reader.Char()
	if !('0' <= nextDigit && nextDigit <= '9') {
		lexer.reader.Redo()
		return base[:len(base)-1], Literal, Integer, nil
	}
	lexer.reader.Next()
	result := append(base, nextDigit)
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if ('0' <= nextDigit && nextDigit <= '9') || nextDigit == '_' {
			result = append(result, nextDigit)
		} else if (nextDigit == 'e') || (nextDigit == 'E') {
			lexer.reader.Next()
			return lexer.tokenizeScientificFloat(append(result, nextDigit))
		} else {
			break
		}
	}
	return result, Literal, Float, nil
}
