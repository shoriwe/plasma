package lexer

import (
	"errors"
	"fmt"
)

type Lexer struct {
	cursor           int
	line             int
	column           int
	sourceCode       string
	sourceCodeLength int
	complete         bool
	peekToken        *Token
}

func (lexer *Lexer) HasNext() bool {
	return !lexer.complete
}

func (lexer *Lexer) tokenizeStringLikeExpressions(stringOpener string, target rune) (string, rune, error) {
	content := stringOpener
	var kind rune = Unknown
	var errorMessage error
	escaped := false
	finish := false
	for ; (lexer.cursor < lexer.sourceCodeLength) && !finish; lexer.cursor++ {
		char := string(lexer.sourceCode[lexer.cursor])
		if escaped {
			switch char {
			case "\\", "'", "\"", "`", "a", "b", "e", "f", "n", "r", "t", "?", "u", "x":
				escaped = false
			default:
				errorMessage = errors.New(fmt.Sprintf("wrong escape at index %d, could not completly define %s", lexer.cursor, content))
				finish = true
			}
		} else {
			switch char {
			case "\n":
				lexer.line++
				lexer.column = 1
			case stringOpener:
				kind = target
				finish = true
			case "\\":
				escaped = true
			}
		}
		content += char
		lexer.column++
	}
	if kind != target {
		errorMessage = errors.New(fmt.Sprintf("No closing at index: %d with value %s", lexer.cursor, content))
	}
	return content, kind, errorMessage
}

func (lexer *Lexer) tokenizeSpecial(target rune) (string, rune, error) {
	stringOpener := string(lexer.sourceCode[lexer.cursor])
	var stringCloser string
	switch stringOpener {
	case "(":
		stringCloser = ")"
	case "[":
		stringCloser = "]"
	case "{":
		stringCloser = "}"
	default:
		stringCloser = stringOpener
	}
	lexer.cursor++
	content := stringOpener
	var kind rune = Unknown
	escaped := false
	finish := false
	var errorMessage error
	for ; (lexer.cursor < lexer.sourceCodeLength) && !finish; lexer.cursor++ {
		char := string(lexer.sourceCode[lexer.cursor])
		if escaped {
			switch char {
			case "\\", "'", "\"", "`", "a", "b", "e", "f", "n", "r", "t", "?", "u", "x", stringOpener:
				escaped = false
			default:
				errorMessage = errors.New(fmt.Sprintf("wrong escape at index %d, could not completly define %s", lexer.cursor, content))
				finish = true
			}
		} else {
			switch char {
			case "\n":
				lexer.line++
				lexer.column = 1
				if stringCloser == "\n" {
					kind = target
					finish = true
				}
			case stringCloser:
				kind = target
				finish = true
			case "\\":
				escaped = true
			}
		}
		content += char
		lexer.column++
	}
	if kind != target {
		errorMessage = errors.New(fmt.Sprintf("No closing at index: %d with value %s", lexer.cursor, content))
	}
	return content, kind, errorMessage
}

func (lexer *Lexer) tokenizeModulusExpression() (string, rune, error) {
	content := "%"
	var computedContent string
	var kind rune = Unknown
	var errorMessage error
	if lexer.cursor+2 >= lexer.sourceCodeLength {
		errorMessage = errors.New(fmt.Sprintf("invalid modulus expression at index: %d", lexer.cursor))
	} else {
		typeOfModulusExpression := string(lexer.sourceCode[lexer.cursor])
		lexer.cursor++
		switch typeOfModulusExpression {
		case "q", "Q": // String2
			content += typeOfModulusExpression
			computedContent, kind, errorMessage = lexer.tokenizeSpecial(String2)
		case "x", "X": // CommandOutput2
			content += typeOfModulusExpression
			computedContent, kind, errorMessage = lexer.tokenizeSpecial(CommandOutput2)
		case "r": // Regexp2
			break
		case "w": // Array 2
			break
		default:
			if isPunctuationPattern.MatchString(typeOfModulusExpression) { // String2
				lexer.cursor--
				content += typeOfModulusExpression
				computedContent, kind, errorMessage = lexer.tokenizeSpecial(String2)
			} else {
				errorMessage = errors.New(fmt.Sprintf("Unknown special pattern char at index %d", lexer.cursor))
			}
		}
	}
	content += computedContent
	return content, kind, errorMessage
}

func (lexer *Lexer) Next() (*Token, error) {
	if lexer.peekToken != nil {
		result := lexer.peekToken
		lexer.peekToken = nil
		return result, nil
	}
	if lexer.cursor == lexer.sourceCodeLength {
		lexer.complete = true
		return &Token{
			String: "EOF",
			Kind:   EOF,
			Line:   lexer.line,
			Column: lexer.column,
			Index:  [2]int{lexer.cursor, -1},
		}, nil
	}
	var tokenizingError error
	var kind rune
	var content string
	column := lexer.column
	line := lexer.line
	index := [2]int{lexer.cursor, -1}
	char := string(lexer.sourceCode[lexer.cursor])
	lexer.cursor++
	switch char {
	case "\n", ";":
		lexer.line++
		lexer.column = 1
		content = char
		kind = Separator
	case "'", "\"": // String1
		content, kind, tokenizingError = lexer.tokenizeStringLikeExpressions(char, String1)
	case "%": // String2 or CommandOutput2 or Regex2 or Array2-Open
		content, kind, tokenizingError = lexer.tokenizeModulusExpression()
	case "`":
		content, kind, tokenizingError = lexer.tokenizeStringLikeExpressions(char, CommandOutput1)
	default:
		panic("Unknown Token char")
	}
	index[1] = lexer.cursor
	return &Token{
		String: content,
		Kind:   kind,
		Line:   line,
		Column: column,
		Index:  index,
	}, tokenizingError
}

func (lexer *Lexer) Peek() (*Token, error) {
	if lexer.peekToken != nil {
		return lexer.peekToken, nil
	}
	var nextError error
	lexer.peekToken, nextError = lexer.Next()
	if nextError != nil {
		return nil, nextError
	}
	return lexer.peekToken, nil
}

func NewLexer(sourceCode string) *Lexer {
	return &Lexer{0, 1, 1, sourceCode, len(sourceCode), false, nil}
}
