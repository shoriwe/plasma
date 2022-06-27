package lexer

import "github.com/shoriwe/gplasma/pkg/errors"

func (lexer *Lexer) tokenizeBinary() *errors.Error {
	if !lexer.reader.HasNext() {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = InvalidDirectValue
		return errors.NewUnknownTokenKindError(lexer.line)
	}
	nextDigit := lexer.reader.Char()
	if !(nextDigit == '0' || nextDigit == '1') {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = InvalidDirectValue
		return errors.NewUnknownTokenKindError(lexer.line)
	}
	lexer.reader.Next()
	lexer.currentToken.append(nextDigit)
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if !(nextDigit == '0' || nextDigit == '1') && nextDigit != '_' {
			break
		}
		lexer.currentToken.append(nextDigit)
	}
	lexer.currentToken.Kind = Literal
	lexer.currentToken.DirectValue = BinaryInteger
	return nil
}
