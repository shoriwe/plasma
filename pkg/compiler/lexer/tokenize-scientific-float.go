package lexer

import "github.com/shoriwe/gplasma/pkg/errors"

func (lexer *Lexer) tokenizeScientificFloat() *errors.Error {
	if !lexer.reader.HasNext() {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = InvalidDirectValue
		return errors.NewUnknownTokenKindError(lexer.line)
	}
	direction := lexer.reader.Char()
	if (direction != '-') && (direction != '+') {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = InvalidDirectValue
		return errors.NewUnknownTokenKindError(lexer.line)
	}
	lexer.reader.Next()
	// Ensure next is a number
	if !lexer.reader.HasNext() {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = InvalidDirectValue
		return errors.NewUnknownTokenKindError(lexer.line)
	}
	lexer.currentToken.append(direction)
	nextDigit := lexer.reader.Char()
	if !('0' <= nextDigit && nextDigit <= '9') {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = InvalidDirectValue
		return errors.NewUnknownTokenKindError(lexer.line)
	}
	lexer.currentToken.append(nextDigit)
	lexer.reader.Next()
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if !('0' <= nextDigit && nextDigit <= '9') && nextDigit != '_' {
			break
		}
		lexer.currentToken.append(nextDigit)
	}
	lexer.currentToken.Kind = Literal
	lexer.currentToken.DirectValue = ScientificFloat
	return nil
}
