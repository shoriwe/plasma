package lexer

import "github.com/shoriwe/gplasma/pkg/errors"

func (lexer *Lexer) tokenizeHexadecimal() *errors.Error {
	if !lexer.reader.HasNext() {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = InvalidDirectValue
		return errors.NewUnknownTokenKindError(lexer.line)
	}
	nextDigit := lexer.reader.Char()
	if !(('0' <= nextDigit && nextDigit <= '9') ||
		('a' <= nextDigit && nextDigit <= 'f') ||
		('A' <= nextDigit && nextDigit <= 'F')) {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = InvalidDirectValue
		return errors.NewUnknownTokenKindError(lexer.line)
	}
	lexer.reader.Next()
	lexer.currentToken.append(nextDigit)
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if !(('0' <= nextDigit && nextDigit <= '9') ||
			('a' <= nextDigit && nextDigit <= 'f') ||
			('A' <= nextDigit && nextDigit <= 'F')) &&
			nextDigit != '_' {
			break
		}
		lexer.currentToken.append(nextDigit)
	}
	lexer.currentToken.Kind = Literal
	lexer.currentToken.DirectValue = HexadecimalInteger
	return nil
}
