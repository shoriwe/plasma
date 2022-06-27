package lexer

func (lexer *Lexer) tokenizeRepeatableOperator(char rune,
	singleDirectValue DirectValue, singleKind Kind,
	doubleDirectValue DirectValue, doubleKind Kind,
	assignSingleDirectValue DirectValue, assignSingleKind Kind,
	assignDoubleDirectValue DirectValue, assignDoubleKind Kind,
) ([]rune, Kind, DirectValue) {
	content := []rune{char}
	kind := singleKind
	directValue := singleDirectValue
	if lexer.reader.HasNext() {
		nextChar := lexer.reader.Char()
		if nextChar == char {
			content = append(content, nextChar)
			lexer.reader.Next()
			kind = doubleKind
			directValue = doubleDirectValue
			if lexer.reader.HasNext() {
				nextNextChar := lexer.reader.Char()
				if nextNextChar == '=' {
					content = append(content, nextNextChar)
					kind = assignDoubleKind
					lexer.reader.Next()
					directValue = assignDoubleDirectValue
				}
			}
		} else if nextChar == '=' {
			kind = assignSingleKind
			content = append(content, nextChar)
			lexer.reader.Next()
			directValue = assignSingleDirectValue
		}
	}
	return content, kind, directValue
}
