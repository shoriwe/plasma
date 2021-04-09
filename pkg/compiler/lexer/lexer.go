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

func (lexer *Lexer) tokenizeStringLikeExpressions(stringOpener string) (string, int, error) {
	content := stringOpener
	var target int
	switch stringOpener {
	case "'":
		target = SingleQuoteString
	case "\"":
		target = DoubleQuoteString
	case "`":
		target = CommandOutput
	}
	var kind = Unknown
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

func (lexer *Lexer) tokenizeNumeric(firstDigit string) (string, int, error) {
	content := firstDigit
	var kind = Integer
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

func guessKind(buffer string) (int, int) {
	switch buffer {
	case PassString:
		return Keyboard, Pass
	case SuperString:
		return Keyboard, Super
	case EndString:
		return Keyboard, End
	case IfString:
		return Keyboard, If
	case ElseString:
		return Keyboard, Else
	case ElifString:
		return Keyboard, Elif
	case WhileString:
		return Keyboard, While
	case ForString:
		return Keyboard, For
	case UntilString:
		return Keyboard, Until
	case SwitchString:
		return Keyboard, Switch
	case CaseString:
		return Keyboard, Case
	case YieldString:
		return Keyboard, Yield
	case ReturnString:
		return Keyboard, Return
	case RetryString:
		return Keyboard, Retry
	case BreakString:
		return Keyboard, Break
	case RedoString:
		return Keyboard, Redo
	case ModuleString:
		return Keyboard, Module
	case DefString:
		return Keyboard, Def
	case LambdaString:
		return Keyboard, Lambda
	case StructString:
		return Keyboard, Struct
	case InterfaceString:
		return Keyboard, Interface
	case GoString:
		return Keyboard, Go
	case ClassString:
		return Keyboard, Class
	case TryString:
		return Keyboard, Try
	case ExceptString:
		return Keyboard, Except
	case FinallyString:
		return Keyboard, Finally
	case AndString:
		return Keyboard, And
	case OrString:
		return Keyboard, Or
	case XorString:
		return Keyboard, Xor
	case InString:
		return Keyboard, In
	case IsInstanceOfString:
		return Keyboard, IsInstanceOf
	case AsyncString:
		return Keyboard, Async
	case AwaitString:
		return Keyboard, Await
	case BEGINString:
		return Keyboard, BEGIN
	case ENDString:
		return Keyboard, END
	case EnumString:
		return Keyboard, Enum
	case NotString:
		return Keyboard, Not
	}
	return IdentifierKind, -1
}

func (lexer *Lexer) tokenizeChars(startingChar string) (string, int, int, error) {
	content := startingChar
	for ; lexer.cursor < lexer.sourceCodeLength; lexer.cursor++ {
		char := string(lexer.sourceCode[lexer.cursor])
		if isNameChar.MatchString(char) {
			content += char
		} else {
			break
		}
	}
	kind, directValue := guessKind(content)
	return content, kind, directValue, nil
}

func (lexer *Lexer) tokenizeComment() (string, int, error) {
	content := ""
	for ; lexer.cursor < lexer.sourceCodeLength; lexer.cursor++ {
		char := string(lexer.sourceCode[lexer.cursor])
		if char == "\n" {
			break
		}
		content += char
	}
	return content, Comment, nil
}

func (lexer *Lexer) tokenizeRepeatableOperator(char string, singleDirectValue int, doubleDirectValue int, assignSingleDirectValue int, assignDoubleDirectValue int) (string, int, int) {
	content := char
	kind := Operator
	directValue := singleDirectValue
	if lexer.cursor < lexer.sourceCodeLength {
		nextChar := string(lexer.sourceCode[lexer.cursor])
		if nextChar == char {
			content += nextChar
			lexer.cursor++
			directValue = doubleDirectValue
			if lexer.cursor < lexer.sourceCodeLength {
				nextNextChar := string(lexer.sourceCode[lexer.cursor])
				if nextNextChar == "=" {
					content += nextNextChar
					kind = Assignment
					lexer.cursor++
					directValue = assignDoubleDirectValue
				}
			}
		} else if nextChar == "=" {
			kind = Assignment
			content += nextChar
			lexer.cursor++
			directValue = assignSingleDirectValue
		}
	}
	return content, kind, directValue
}

func (lexer *Lexer) tokenizeNotRepeatableOperator(char string, single int, assign int) (string, int, int) {
	content := char
	kind := Operator
	directValue := single
	if lexer.cursor < lexer.sourceCodeLength {
		nextChar := string(lexer.sourceCode[lexer.cursor])
		if nextChar == "=" {
			kind = Assignment
			directValue = assign
			content += nextChar
			lexer.cursor++
		}
	}
	return content, kind, directValue
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
	var kind int
	var directValue int
	var content string
	line := lexer.line
	index := lexer.cursor
	char := string(lexer.sourceCode[lexer.cursor])
	lexer.cursor++
	switch char {
	case NewLineString:
		lexer.line++
		content = char
		kind = Separator
		directValue = NewLine
	case SemiColonString:
		content = char
		kind = Separator
		directValue = SemiColon
	case ColonString:
		content = char
		directValue = Colon
		kind = Punctuation
	case CommaString:
		content = char
		directValue = Comma
		kind = Punctuation
	case OpenParenthesesString:
		content = char
		directValue = OpenParentheses
		kind = Punctuation
	case CloseParenthesesString:
		content = char
		directValue = CloseParentheses
		kind = Punctuation
	case OpenSquareBracketString:
		content = char
		directValue = OpenSquareBracket
		kind = Punctuation
	case CloseSquareBracketString:
		content = char
		directValue = CloseSquareBracket
		kind = Punctuation
	case OpenBraceString:
		content = char
		directValue = OpenBrace
		kind = Punctuation
	case CloseBraceString:
		content = char
		directValue = CloseBrace
		kind = Punctuation
	case DollarSignString:
		content = char
		directValue = DollarSign
		kind = Punctuation
	case DotString:
		content = char
		directValue = Dot
		kind = Punctuation
	case WhiteSpaceString:
		directValue = Whitespace
		kind = Whitespace
		content = char
	case TabString:
		directValue = Tab
		kind = Whitespace
		content = char
	case CommentString:
		content, kind, tokenizingError = lexer.tokenizeComment()
		content = "#" + content
	case "'", "\"": // String1
		content, kind, tokenizingError = lexer.tokenizeStringLikeExpressions(char)
	case "`":
		content, kind, tokenizingError = lexer.tokenizeStringLikeExpressions(char)
	case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
		content, kind, tokenizingError = lexer.tokenizeNumeric(char)
	case StarString:
		content, kind, directValue = lexer.tokenizeRepeatableOperator(char, Star, PowerOf, StarAssign, PowerOfAssign)
	case DivString:
		content, kind, directValue = lexer.tokenizeRepeatableOperator(char, Div, FloorDiv, DivAssign, FloorDivAssign)
	case LessThanString:
		content, kind, directValue = lexer.tokenizeRepeatableOperator(char, LessThan, BitwiseLeft, LessOrEqualThan, BitwiseLeftAssign)
	case GreatThanString:
		content, kind, directValue = lexer.tokenizeRepeatableOperator(char, GreaterThan, BitwiseRight, GreaterOrEqualThan, BitwiseRightAssign)
	case AddString:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, Add, AddAssign)
	case SubString:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, Sub, SubAssign)
	case ModulusString:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, Modulus, ModulusAssign)
	case BitwiseXorString:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, BitwiseXor, BitwiseXorAssign)
	case BitWiseAndString:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, BitWiseAnd, BitWiseAndAssign)
	case BitwiseOrString:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, BitwiseOr, BitwiseOrAssign)
	case SignNotString:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, SignNot, NotEqual)
	case NegateBitsString:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, NegateBits, NegateBitsAssign)
	case EqualsString:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, Assign, Equals)
	case "\\":
		content = char
		if lexer.cursor < lexer.sourceCodeLength {
			nextChar := string(lexer.sourceCode[lexer.cursor])
			if nextChar != "\n" {
				return nil, errors.New("pending line escape not followed by a new line")
			}
			content += "\n"
			lexer.cursor++
		}
		kind = PendingEscape
	default:
		if char == "b" {
			if lexer.cursor < lexer.sourceCodeLength {
				nextChar := string(lexer.sourceCode[lexer.cursor])
				if nextChar == "'" || nextChar == "\"" {
					var byteStringPart string
					lexer.cursor++
					byteStringPart, kind, tokenizingError = lexer.tokenizeStringLikeExpressions(nextChar)
					content = char + byteStringPart
					if kind != Unknown {
						kind = ByteString
					}
					break
				}
			}
		}
		content, kind, directValue, tokenizingError = lexer.tokenizeChars(char)
	}
	return &Token{
		DirectValue: directValue,
		String:      content,
		Kind:        kind,
		Line:        line,
		Index:       index,
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
