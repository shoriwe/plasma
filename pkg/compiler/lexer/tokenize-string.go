package lexer

import "github.com/shoriwe/gplasma/pkg/errors"

func (lexer *Lexer) tokenizeStringLikeExpressions(stringOpener rune) ([]rune, Kind, DirectValue, *errors.Error) {
	content := []rune{stringOpener}
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
	var tokenizingError *errors.Error
	escaped := false
	finish := false
	for ; lexer.reader.HasNext() && !finish; lexer.reader.Next() {
		char := lexer.reader.Char()
		if escaped {
			switch char {
			case '\\', '\'', '"', '`', 'a', 'b', 'e', 'f', 'n', 'r', 't', '?', 'u', 'x':
				escaped = false
			default:
				tokenizingError = errors.New(lexer.line, "invalid escape sequence", errors.LexingError)
				finish = true
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
		content = append(content, char)
	}
	if directValue != target {
		tokenizingError = errors.New(lexer.line, "string never closed", errors.LexingError)
	}
	return content, Literal, directValue, tokenizingError
}
