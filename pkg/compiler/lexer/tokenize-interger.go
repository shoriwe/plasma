package lexer

import "github.com/shoriwe/gplasma/pkg/errors"

func (lexer *Lexer) tokenizeInteger() *errors.Error {
	if !lexer.reader.HasNext() {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = Integer
		return nil
	}
	nextDigit := lexer.reader.Char()
	if nextDigit == '.' {
		lexer.reader.Next()
		lexer.currentToken.append(nextDigit)
		return lexer.tokenizeFloat()
	} else if nextDigit == 'e' || nextDigit == 'E' {
		lexer.reader.Next()
		lexer.currentToken.append(nextDigit)
		return lexer.tokenizeScientificFloat()
	} else if !('0' <= nextDigit && nextDigit <= '9') {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = Integer
		return nil
	}
	lexer.reader.Next()
	lexer.currentToken.append(nextDigit)
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if nextDigit == 'e' || nextDigit == 'E' {
			lexer.currentToken.append(nextDigit)
			return lexer.tokenizeScientificFloat()
		} else if nextDigit == '.' {
			lexer.reader.Next()
			lexer.currentToken.append(nextDigit)
			return lexer.tokenizeFloat()
		} else if ('0' <= nextDigit && nextDigit <= '9') || nextDigit == '_' {
			lexer.currentToken.append(nextDigit)
		} else {
			break
		}
	}
	lexer.currentToken.Kind = Literal
	lexer.currentToken.DirectValue = Integer
	return nil
}
