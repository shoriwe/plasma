package lexer

import "github.com/shoriwe/gplasma/pkg/errors"

func (lexer *Lexer) tokenizeInteger(base []rune) ([]rune, Kind, DirectValue, *errors.Error) {
	if !lexer.reader.HasNext() {
		return base, Literal, Integer, nil
	}
	nextDigit := lexer.reader.Char()
	if nextDigit == '.' {
		lexer.reader.Next()
		return lexer.tokenizeFloat(append(base, nextDigit))
	} else if nextDigit == 'e' || nextDigit == 'E' {
		lexer.reader.Next()
		return lexer.tokenizeScientificFloat(append(base, nextDigit))
	} else if !('0' <= nextDigit && nextDigit <= '9') {
		return base, Literal, Integer, nil
	}
	lexer.reader.Next()
	result := append(base, nextDigit)
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if nextDigit == 'e' || nextDigit == 'E' {
			return lexer.tokenizeScientificFloat(append(result, nextDigit))
		} else if nextDigit == '.' {
			lexer.reader.Next()
			return lexer.tokenizeFloat(append(result, nextDigit))
		} else if ('0' <= nextDigit && nextDigit <= '9') || nextDigit == '_' {
			result = append(result, nextDigit)
		} else {
			break
		}
	}
	return result, Literal, Integer, nil
}
