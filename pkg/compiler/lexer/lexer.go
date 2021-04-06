package lexer

import (
	"errors"
	"fmt"
	"regexp"
)

type Lexer struct {
	cursor           int
	lastToken        *Token
	line             int
	sourceCode       string
	sourceCodeLength int
	complete         bool
	peekToken        *Token
}

func (lexer *Lexer) HasNext() bool {
	return !lexer.complete
}

func (lexer *Lexer) tokenizeStringLikeExpressions(stringOpener string) (string, rune, error) {
	content := stringOpener
	var target rune
	switch stringOpener {
	case "'":
		target = SingleQuoteString
	case "\"":
		target = DoubleQuoteString
	case "`":
		target = CommandOutput
	}
	var kind rune = Unknown
	var tokenizingError error
	escaped := false
	finish := false
	for ; (lexer.cursor < lexer.sourceCodeLength) && !finish; lexer.cursor++ {
		char := string(lexer.sourceCode[lexer.cursor])
		if escaped {
			switch char {
			case "\\", "'", "\"", "`", "a", "b", "e", "f", "n", "r", "t", "?", "u", "x":
				escaped = false
			default:
				tokenizingError = errors.New(fmt.Sprintf("wrong escape at index %d, could not completly define %s", lexer.cursor, content))
				finish = true
			}
		} else {
			switch char {
			case "\n":
				lexer.line++
			case stringOpener:
				kind = target
				finish = true
			case "\\":
				escaped = true
			}
		}
		content += char
	}
	if kind != target {
		tokenizingError = errors.New(fmt.Sprintf("No closing at index: %d with value %s", lexer.cursor, content))
	}
	return content, kind, tokenizingError
}

func (lexer *Lexer) tokenizeNumeric(firstDigit string) (string, rune, error) {
	content := firstDigit
	var kind rune = Integer
	var tokenizingError error
tokenizingLoop:
	for ; lexer.cursor < lexer.sourceCodeLength; lexer.cursor++ {
		char := string(lexer.sourceCode[lexer.cursor])
		switch char {
		case "_", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "a", "A", "c", "C", "d", "D", "f", "F":
			switch kind {
			case Integer | Float | ScientificFloat:
				switch char {
				case "_", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
					content += char
				default:
					tokenizingError = errors.New(fmt.Sprintf("invalid numeric declaration at index: %d", lexer.cursor))
					break tokenizingLoop
				}
			case HexadecimalInteger:
				content += char
			case BinaryInteger:
				switch char {
				case "_", "1", "0":
					content += char
				default:
					tokenizingError = errors.New(fmt.Sprintf("invalid numeric declaration at index: %d", lexer.cursor))
					break tokenizingLoop
				}
			case OctalInteger:
				switch char {
				case "_", "1", "2", "3", "4", "5", "6", "7", "0":
					content += char
				default:
					tokenizingError = errors.New(fmt.Sprintf("invalid numeric declaration at index: %d", lexer.cursor))
					break tokenizingLoop
				}
			}
		case ".":
			if kind == Float || kind == ScientificFloat {
				break tokenizingLoop
			}
			kind = Float
			content += "."
		case "x", "X":
			if content != "0" {
				tokenizingError = errors.New(fmt.Sprintf("invalid numeric declaration at index: %d", lexer.cursor))
				break tokenizingLoop
			}
			if lexer.cursor+1 < lexer.sourceCodeLength {
				nextChar := string(lexer.sourceCode[lexer.cursor+1])
				switch nextChar {
				case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "a", "A", "b", "B", "c", "C", "d", "D", "e", "E", "f", "F":
					kind = HexadecimalInteger
					content += "x"
					lexer.cursor++
				default:
					tokenizingError = errors.New(fmt.Sprintf("invalid numeric declaration at index: %d", lexer.cursor))
					break tokenizingLoop
				}
			}
		case "b", "B":
			if kind == HexadecimalInteger {
				content += "b"
				continue
			}
			if content != "0" {
				tokenizingError = errors.New(fmt.Sprintf("invalid numeric declaration at index: %d", lexer.cursor))
				break tokenizingLoop
			}
			if lexer.cursor+1 < lexer.sourceCodeLength {
				nextChar := string(lexer.sourceCode[lexer.cursor+1])
				switch nextChar {
				case "1", "0":
					kind = BinaryInteger
					content += "b"
					lexer.cursor++
				default:
					tokenizingError = errors.New(fmt.Sprintf("invalid numeric declaration at index: %d", lexer.cursor))
					break tokenizingLoop
				}
			}
		case "o", "O":
			if content != "0" {
				tokenizingError = errors.New(fmt.Sprintf("invalid numeric declaration at index: %d", lexer.cursor))
				break tokenizingLoop
			}
			if lexer.cursor+1 < lexer.sourceCodeLength {
				nextChar := string(lexer.sourceCode[lexer.cursor+1])
				switch nextChar {
				case "1", "2", "3", "4", "5", "6", "7", "0":
					kind = OctalInteger
					content += "o"
					lexer.cursor++
				default:
					tokenizingError = errors.New(fmt.Sprintf("invalid numeric declaration at index: %d", lexer.cursor))
					break tokenizingLoop
				}
			}
		case "e", "E":
			if kind == HexadecimalInteger {
				content += "e"
				continue
			}
			if kind != Float && kind != Integer {
				tokenizingError = errors.New(fmt.Sprintf("Multiple scientific syntax in the same number at index: %d", lexer.cursor))
				break tokenizingLoop
			}
			if lexer.cursor+2 < lexer.sourceCodeLength {
				nextChar := string(lexer.sourceCode[lexer.cursor+1])
				nextNextChar := string(lexer.sourceCode[lexer.cursor+2])
				if nextChar == "-" {
					switch nextNextChar {
					case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
						content += "e-" + nextNextChar
						kind = ScientificFloat
						lexer.cursor += 2
						continue
					}
				}
				tokenizingError = errors.New(fmt.Sprintf("Invalid scientific number declaration at index: %d", lexer.cursor))
				break tokenizingLoop
			}
			tokenizingError = errors.New(fmt.Sprintf("Incoplete Scientific number declaration at index: %d", lexer.cursor))
			break tokenizingLoop

		default:
			break tokenizingLoop
		}
	}
	if kind == Float && content[len(content)-1] == '.' {
		content = content[:len(content)-1]
		kind = Integer
		lexer.cursor--
	}
	return content, kind, tokenizingError
}

var isNameChar = regexp.MustCompile("[_a-zA-Z0-9]")
var isConstant = regexp.MustCompile("[A-Z]+[_a-zA-Z0-9]*")

func GuessKind(buffer string) rune {
	switch buffer {
	case Super, End, If, Else, Elif, While, For, Until, Switch, Case, Yield, Return, Retry, Break, Redo, Module, Def, Lambda, Struct, Interface, Go, Class, Try, Except, Finally, And, Or, Xor, In, IsInstanceOf, When, Async, Await, BEGIN, END, Enum:
		return Keyboard
	}
	if isConstant.MatchString(buffer) {
		return ConstantKind
	}
	return IdentifierKind
}

func (lexer *Lexer) tokenizeChars(startingChar string) (string, rune, error) {
	content := startingChar
	for ; lexer.cursor < lexer.sourceCodeLength; lexer.cursor++ {
		char := string(lexer.sourceCode[lexer.cursor])
		if isNameChar.MatchString(char) {
			content += char
		} else {
			break
		}
	}
	return content, GuessKind(content), nil
}

func (lexer *Lexer) next() (*Token, error) {
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
			Index:  lexer.cursor,
		}, nil
	}
	var tokenizingError error
	var kind rune
	var content string
	line := lexer.line
	index := lexer.cursor
	char := string(lexer.sourceCode[lexer.cursor])
	lexer.cursor++
	switch char {
	case "\n", ";":
		lexer.line++
		content = char
		kind = Separator
	case "'", "\"": // String1
		content, kind, tokenizingError = lexer.tokenizeStringLikeExpressions(char)
	case "`":
		content, kind, tokenizingError = lexer.tokenizeStringLikeExpressions(char)
	case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
		content, kind, tokenizingError = lexer.tokenizeNumeric(char)
	case ":", ",", "(", ")", "[", "]", "{", "}", "$", ".":
		lexer.line++
		content = char
		kind = Punctuation
	case "*", "/":
		content = char
		kind = Operator
		if lexer.cursor < lexer.sourceCodeLength {
			nextChar := string(lexer.sourceCode[lexer.cursor])
			if nextChar == char {
				content += nextChar
				lexer.cursor++
				if lexer.cursor < lexer.sourceCodeLength {
					nextNextChar := string(lexer.sourceCode[lexer.cursor])
					if nextNextChar == "=" {
						content += nextNextChar
						kind = Assignment
						lexer.cursor++
					}
				}
			} else if nextChar == "=" {
				kind = Assignment
				content += nextChar
				lexer.cursor++
			}
		}
	case "+", "-", "%", "^", "&", "|", "!", "~":
		content = char
		kind = Operator
		if lexer.cursor < lexer.sourceCodeLength {
			nextChar := string(lexer.sourceCode[lexer.cursor])
			if nextChar == "=" {
				kind = Assignment
				content += nextChar
				lexer.cursor++
			}
		}
	case "<", ">":
		content = char
		kind = Comparator
		if lexer.cursor < lexer.sourceCodeLength {
			nextChar := string(lexer.sourceCode[lexer.cursor])
			if nextChar == char {
				content += nextChar
				kind = Operator
				lexer.cursor++
				if lexer.cursor < lexer.sourceCodeLength {
					nextNextChar := string(lexer.sourceCode[lexer.cursor])
					if nextNextChar == "=" {
						content += nextNextChar
						kind = Assignment
						lexer.cursor++
					}
				}
			} else if nextChar == "=" {
				content += nextChar
				lexer.cursor++
			}
		}
	case "=":
		content += char
		kind = Assignment
		if lexer.cursor+1 < lexer.sourceCodeLength {
			nextChar := string(lexer.sourceCode[lexer.cursor+1])
			if nextChar == "=" {
				kind = Comparator
				content += nextChar
				lexer.cursor++
			}
		}
	case " ", "\t":
		kind = Whitespace
		content = char
	default:
		content, kind, tokenizingError = lexer.tokenizeChars(char)
	}
	return &Token{
		String: content,
		Kind:   kind,
		Line:   line,
		Index:  index,
	}, tokenizingError
}

/*
	This function will yield just the necessary token, this means not repeated separators
*/
func (lexer *Lexer) Next() (*Token, error) {
nextTokenLoop:
	for ; lexer.HasNext(); {
		token, tokenizationError := lexer.next()
		if tokenizationError != nil {
			return nil, tokenizationError
		}
		switch token.Kind {
		case Whitespace:
			continue
		case Separator:
			if lexer.lastToken == nil {
				continue
			}
			switch lexer.lastToken.Kind {
			case Operator:
				continue
			case Comparator:
				continue
			case Separator:
				continue
			default:
				lexer.lastToken = token
				break nextTokenLoop
			}
		default:
			lexer.lastToken = token
			break nextTokenLoop
		}
	}
	return lexer.lastToken, nil
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
	return &Lexer{0, nil, 1, sourceCode, len(sourceCode), false, nil}
}
