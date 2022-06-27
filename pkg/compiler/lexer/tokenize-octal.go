package lexer

import "github.com/shoriwe/gplasma/pkg/errors"

func (lexer *Lexer) tokenizeOctal() *errors.Error {
	if !lexer.reader.HasNext() {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = InvalidDirectValue
		return errors.NewUnknownTokenKindError(lexer.line)
	}
	nextDigit := lexer.reader.Char()
	if !('0' <= nextDigit && nextDigit <= '7') {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = InvalidDirectValue
		return errors.NewUnknownTokenKindError(lexer.line)
	}
	lexer.reader.Next()
	lexer.currentToken.append(nextDigit)
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if !('0' <= nextDigit && nextDigit <= '7') && nextDigit != '_' {
			break
		}
		lexer.currentToken.append(nextDigit)
	}
	lexer.currentToken.Kind = Literal
	lexer.currentToken.DirectValue = OctalInteger
	return nil
}
