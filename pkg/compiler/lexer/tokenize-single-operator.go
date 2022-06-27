package lexer

func (lexer *Lexer) tokenizeSingleOperator(char rune,
	single DirectValue, singleKind Kind,
	assign DirectValue, assignKind Kind) ([]rune, Kind, DirectValue) {
	content := []rune{char}
	kind := singleKind
	directValue := single
	if lexer.reader.HasNext() {
		nextChar := lexer.reader.Char()
		if nextChar == '=' {
			kind = assignKind
			directValue = assign
			content = append(content, nextChar)
			lexer.reader.Next()
		}
	}
	return content, kind, directValue
}
