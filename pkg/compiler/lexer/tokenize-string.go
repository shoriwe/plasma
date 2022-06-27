package lexer

import "github.com/shoriwe/gplasma/pkg/errors"

func (lexer *Lexer) tokenizeStringLikeExpressions(stringOpener rune) *errors.Error {
	var target DirectValue
	switch stringOpener {
	case '\'':
		target = SingleQuoteString
	case '"':
		target = DoubleQuoteString
	case '`':
		target = CommandOutput
	}
	var directValue = InvalidDirectValue
	escaped := false
	finish := false
	for ; lexer.reader.HasNext() && !finish; lexer.reader.Next() {
		char := lexer.reader.Char()
		if escaped {
			switch char {
			case '\\', '\'', '"', '`', 'a', 'b', 'e', 'f', 'n', 'r', 't', '?', 'u', 'x':
				escaped = false
			default:
				return errors.New(lexer.line, "invalid escape sequence", errors.LexingError)
			}
		} else {
			switch char {
			case '\n':
				lexer.line++
			case stringOpener:
				directValue = target
				finish = true
			case '\\':
				escaped = true
			}
		}
		lexer.currentToken.append(char)
	}
	if directValue != target {
		return errors.New(lexer.line, "string never closed", errors.LexingError)
	}
	lexer.currentToken.Kind = Literal
	lexer.currentToken.DirectValue = directValue
	return nil
}
