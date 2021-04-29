package lexer

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	hexDigitPattern    = regexp.MustCompile("[0-9a-fA-F]")
	binaryDigitPattern = regexp.MustCompile("[01]")
	octalDigitPattern  = regexp.MustCompile("[0-7]")
	digitPattern       = regexp.MustCompile("[0-9]")
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

func (lexer *Lexer) tokenizeStringLikeExpressions(stringOpener string) (string, int, int, error) {
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
	var directValue = Unknown
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
				directValue = target
				finish = true
			case "\\":
				escaped = true
			}
		}
		content += char
	}
	if directValue != target {
		tokenizingError = errors.New(fmt.Sprintf("No closing at index: %d with value %s", lexer.cursor, content))
	}
	return content, Literal, directValue, tokenizingError
}

func (lexer *Lexer) tokenizeHexadecimal(letterX string) (string, int, int, error) {
	result := "0" + letterX
	if !(lexer.cursor < lexer.sourceCodeLength) {
		return result, Literal, Unknown, errors.New(fmt.Sprintf("could not determine the kind of the token %s at line %d", result, lexer.line))
	}
	nextDigit := string(lexer.sourceCode[lexer.cursor])
	if !hexDigitPattern.MatchString(nextDigit) {
		return result, Literal, Unknown, errors.New(fmt.Sprintf("could not determine the kind of the token %s at line %d", result, lexer.line))
	}
	lexer.cursor++
	result += nextDigit
	for ; lexer.cursor < lexer.sourceCodeLength; lexer.cursor++ {
		nextDigit = string(lexer.sourceCode[lexer.cursor])
		if !hexDigitPattern.MatchString(nextDigit) && nextDigit != "_" {
			return result, Literal, HexadecimalInteger, nil
		}
		result += nextDigit
	}
	return result, Literal, HexadecimalInteger, nil
}

func (lexer *Lexer) tokenizeBinary(letterB string) (string, int, int, error) {
	result := "0" + letterB
	if !(lexer.cursor < lexer.sourceCodeLength) {
		return result, Literal, Unknown, errors.New(fmt.Sprintf("could not determine the kind of the token %s at line %d", result, lexer.line))
	}
	nextDigit := string(lexer.sourceCode[lexer.cursor])
	if !binaryDigitPattern.MatchString(nextDigit) {
		return result, Literal, Unknown, errors.New(fmt.Sprintf("could not determine the kind of the token %s at line %d", result, lexer.line))
	}
	lexer.cursor++
	result += nextDigit
	for ; lexer.cursor < lexer.sourceCodeLength; lexer.cursor++ {
		nextDigit = string(lexer.sourceCode[lexer.cursor])
		if !binaryDigitPattern.MatchString(nextDigit) && nextDigit != "_" {
			return result, Literal, BinaryInteger, nil
		}
		result += nextDigit
	}
	return result, Literal, BinaryInteger, nil
}

func (lexer *Lexer) tokenizeOctal(letterO string) (string, int, int, error) {
	result := "0" + letterO
	if !(lexer.cursor < lexer.sourceCodeLength) {
		return result, Literal, Unknown, errors.New(fmt.Sprintf("could not determine the kind of the token %s at line %d", result, lexer.line))
	}
	nextDigit := string(lexer.sourceCode[lexer.cursor])
	if !octalDigitPattern.MatchString(nextDigit) {
		return result, Literal, Unknown, errors.New(fmt.Sprintf("could not determine the kind of the token %s at line %d", result, lexer.line))
	}
	lexer.cursor++
	result += nextDigit
	for ; lexer.cursor < lexer.sourceCodeLength; lexer.cursor++ {
		nextDigit = string(lexer.sourceCode[lexer.cursor])
		if !octalDigitPattern.MatchString(nextDigit) && nextDigit != "_" {
			return result, Literal, OctalInteger, nil
		}
		result += nextDigit
	}
	return result, Literal, OctalInteger, nil
}

func (lexer *Lexer) tokenizeScientificFloat(base string) (string, int, int, error) {
	result := base
	if !(lexer.cursor < lexer.sourceCodeLength) {
		return result, Literal, Unknown, errors.New(fmt.Sprintf("could not determine the kind of the token %s at line %d", result, lexer.line))
	}
	direction := string(lexer.sourceCode[lexer.cursor])
	if direction != "-" && direction != "+" {
		return result, Literal, Unknown, errors.New(fmt.Sprintf("could not determine the kind of the token %s at line %d", result, lexer.line))
	}
	lexer.cursor++
	// Ensure next is a number
	if !(lexer.cursor < lexer.sourceCodeLength) {
		return result, Literal, Unknown, errors.New(fmt.Sprintf("could not determine the kind of the token %s at line %d", result, lexer.line))
	}
	result += direction
	nextDigit := string(lexer.sourceCode[lexer.cursor])
	if !digitPattern.MatchString(nextDigit) {
		return result, Literal, Unknown, errors.New(fmt.Sprintf("could not determine the kind of the token %s at line %d", result, lexer.line))
	}
	result += nextDigit
	lexer.cursor++
	for ; lexer.cursor < lexer.sourceCodeLength; lexer.cursor++ {
		nextDigit = string(lexer.sourceCode[lexer.cursor])
		if !digitPattern.MatchString(nextDigit) && nextDigit != "_" {
			return result, Literal, ScientificFloat, nil
		}
		result += nextDigit
	}
	return result, Literal, ScientificFloat, nil
}

func (lexer *Lexer) tokenizeFloat(base string) (string, int, int, error) {
	if !(lexer.cursor < lexer.sourceCodeLength) {
		lexer.cursor--
		return base[:len(base)-1], Literal, Integer, nil
	}
	nextDigit := string(lexer.sourceCode[lexer.cursor])
	if !digitPattern.MatchString(nextDigit) {
		lexer.cursor--
		return base[:len(base)-1], Literal, Integer, nil
	}
	lexer.cursor++
	result := base + nextDigit
	for ; lexer.cursor < lexer.sourceCodeLength; lexer.cursor++ {
		nextDigit = string(lexer.sourceCode[lexer.cursor])
		if digitPattern.MatchString(nextDigit) || nextDigit == "_" {
			result += nextDigit
		} else if (nextDigit == "e") || (nextDigit == "E") {
			return lexer.tokenizeScientificFloat(result + nextDigit)
		} else {
			break
		}
	}
	return result, Literal, Float, nil
}

func (lexer *Lexer) tokenizeInteger(base string) (string, int, int, error) {
	if !(lexer.cursor < lexer.sourceCodeLength) {
		return base, Literal, Integer, nil
	}
	nextDigit := string(lexer.sourceCode[lexer.cursor])
	if nextDigit == "." {
		lexer.cursor++
		return lexer.tokenizeFloat(base + nextDigit)
	} else if nextDigit == "e" || nextDigit == "E" {
		lexer.cursor++
		return lexer.tokenizeScientificFloat(base + nextDigit)
	} else if !digitPattern.MatchString(nextDigit) {
		return base, Literal, Integer, nil
	}
	lexer.cursor++
	result := base + nextDigit
	for ; lexer.cursor < lexer.sourceCodeLength; lexer.cursor++ {
		nextDigit = string(lexer.sourceCode[lexer.cursor])
		if nextDigit == "e" || nextDigit == "E" {
			return lexer.tokenizeScientificFloat(result + nextDigit)
		} else if nextDigit == "." {
			return lexer.tokenizeFloat(result + nextDigit)
		} else if digitPattern.MatchString(nextDigit) || nextDigit == "_" {
			result += nextDigit
		} else {
			break
		}
	}
	return result, Literal, Integer, nil
}

func (lexer *Lexer) tokenizeNumeric(firstDigit string) (string, int, int, error) {
	if !(lexer.cursor < lexer.sourceCodeLength) {
		return firstDigit, Literal, Integer, nil
	}
	nextChar := string(lexer.sourceCode[lexer.cursor])
	lexer.cursor++
	switch firstDigit {
	case "0": // In this scenario it can be a float,  scientific float, integer, hex integer, octal integer or binary integer
		switch nextChar {
		case "x", "X": // Hexadecimal
			return lexer.tokenizeHexadecimal(nextChar)
		case "b", "B": // Binary
			return lexer.tokenizeBinary(nextChar)
		case "o", "O": // Octal
			return lexer.tokenizeOctal(nextChar)
		case "e", "E": // Scientific float
			return lexer.tokenizeScientificFloat(firstDigit + nextChar)
		case ".": // Maybe a float
			return lexer.tokenizeFloat(firstDigit + nextChar) // Integer, Float Or Scientific Float
		default:
			if digitPattern.MatchString(nextChar) || nextChar == "_" {
				return lexer.tokenizeInteger(firstDigit + nextChar) // Integer, Float or Scientific Float
			}
		}
	default:
		switch nextChar {
		case "e", "E": // Scientific float
			return lexer.tokenizeScientificFloat(firstDigit + nextChar)
		case ".": // Maybe a float
			return lexer.tokenizeFloat(firstDigit + nextChar) // Integer, Float Or Scientific Float
		default:
			if digitPattern.MatchString(nextChar) || nextChar == "_" {
				return lexer.tokenizeInteger(firstDigit + nextChar) // Integer, Float or Scientific Float
			}
		}
	}
	lexer.cursor--
	return firstDigit, Literal, Integer, nil
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
	case UnlessString:
		return Keyboard, Unless
	case ElseString:
		return Keyboard, Else
	case ElifString:
		return Keyboard, Elif
	case WhileString:
		return Keyboard, While
	case DoString:
		return Keyboard, Do
	case ForString:
		return Keyboard, For
	case UntilString:
		return Keyboard, Until
	case SwitchString:
		return Keyboard, Switch
	case CaseString:
		return Keyboard, Case
	case DefaultString:
		return Keyboard, Default
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
		return Comparator, And
	case OrString:
		return Comparator, Or
	case XorString:
		return Comparator, Xor
	case InString:
		return Comparator, In
	case IsInstanceOfString: // This is a method like super
		return Keyboard, IsInstanceOf
	case AsyncString:
		return Keyboard, Async
	case AsString:
		return Keyboard, As
	case RaiseString:
		return Keyboard, Raise
	case AwaitString:
		return AwaitKeyboard, Await
	case BEGINString:
		return Keyboard, BEGIN
	case ENDString:
		return Keyboard, END
	case EnumString:
		return Keyboard, Enum
	case NotString: // Unary operator
		return Operator, Not
	case TrueString:
		return Boolean, True
	case FalseString:
		return Boolean, False
	case NoneString:
		return NoneType, None
	case DeferString:
		return Keyboard, Defer
	case GoToString:
		return Keyboard, GoTo
	case ContextString:
		return Keyboard, Context
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

func (lexer *Lexer) tokenizeRepeatableOperator(char string, singleDirectValue int, singleKind int, doubleDirectValue int, doubleKind int, assignSingleDirectValue int, assignDoubleDirectValue int, assignSingleKind int, assignDoubleKind int) (string, int, int) {
	content := char
	kind := singleKind
	directValue := singleDirectValue
	if lexer.cursor < lexer.sourceCodeLength {
		nextChar := string(lexer.sourceCode[lexer.cursor])
		if nextChar == char {
			content += nextChar
			lexer.cursor++
			kind = doubleKind
			directValue = doubleDirectValue
			if lexer.cursor < lexer.sourceCodeLength {
				nextNextChar := string(lexer.sourceCode[lexer.cursor])
				if nextNextChar == "=" {
					content += nextNextChar
					kind = assignDoubleKind
					lexer.cursor++
					directValue = assignDoubleDirectValue
				}
			}
		} else if nextChar == "=" {
			kind = assignSingleKind
			content += nextChar
			lexer.cursor++
			directValue = assignSingleDirectValue
		}
	}
	return content, kind, directValue
}

func (lexer *Lexer) tokenizeNotRepeatableOperator(char string, single int, singleKind int, assign int, assignKind int) (string, int, int) {
	content := char
	kind := singleKind
	directValue := single
	if lexer.cursor < lexer.sourceCodeLength {
		nextChar := string(lexer.sourceCode[lexer.cursor])
		if nextChar == "=" {
			kind = assignKind
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
			String:      "EOF",
			DirectValue: EOF,
			Kind:        EOF,
			Line:        lexer.line,
			Index:       lexer.cursor,
		}, nil
	}
	var tokenizingError error
	var kind int
	var content string
	directValue := Unknown
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
		content, kind, directValue, tokenizingError = lexer.tokenizeStringLikeExpressions(char)
	case "`":
		content, kind, directValue, tokenizingError = lexer.tokenizeStringLikeExpressions(char)
	case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
		content, kind, directValue, tokenizingError = lexer.tokenizeNumeric(char)
	case StarString:
		content, kind, directValue = lexer.tokenizeRepeatableOperator(char, Star, Operator, PowerOf, Operator, StarAssign, PowerOfAssign, Assignment, Assignment)
	case DivString:
		content, kind, directValue = lexer.tokenizeRepeatableOperator(char, Div, Operator, FloorDiv, Operator, DivAssign, FloorDivAssign, Assignment, Assignment)
	case LessThanString:
		content, kind, directValue = lexer.tokenizeRepeatableOperator(char, LessThan, Comparator, BitwiseLeft, Operator, LessOrEqualThan, BitwiseLeftAssign, Comparator, Assignment)
	case GreatThanString:
		content, kind, directValue = lexer.tokenizeRepeatableOperator(char, GreaterThan, Comparator, BitwiseRight, Operator, GreaterOrEqualThan, BitwiseRightAssign, Comparator, Assignment)
	case AddString:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, Add, Operator, AddAssign, Assignment)
	case SubString:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, Sub, Operator, SubAssign, Assignment)
	case ModulusString:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, Modulus, Operator, ModulusAssign, Assignment)
	case BitwiseXorString:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, BitwiseXor, Operator, BitwiseXorAssign, Assignment)
	case BitWiseAndString:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, BitWiseAnd, Operator, BitWiseAndAssign, Assignment)
	case BitwiseOrString:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, BitwiseOr, Operator, BitwiseOrAssign, Assignment)
	case SignNotString:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, SignNot, Operator, NotEqual, Comparator)
	case NegateBitsString:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, NegateBits, Operator, NegateBitsAssign, Assignment)
	case EqualsString:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, Assign, Assignment, Equals, Comparator)
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
					byteStringPart, kind, directValue, tokenizingError = lexer.tokenizeStringLikeExpressions(nextChar)
					content = char + byteStringPart
					if directValue != Unknown {
						directValue = ByteString
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
	token, tokenizingError := lexer.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if token.Kind == Comment {
		return lexer.Next()
	}
	if token.Kind == Whitespace {
		return lexer.Next()
	}
	if token.Kind == Separator {
		if lexer.lastToken == nil {
			return lexer.Next()
		}
		switch lexer.lastToken.Kind {
		case Separator:
			return lexer.Next()
		case Operator, Comparator:
			return lexer.Next()
		default:
			break
		}
		switch lexer.lastToken.DirectValue {
		case Comma, OpenParentheses, OpenSquareBracket, OpenBrace:
			return lexer.Next()
		default:
			break
		}
	}
	lexer.lastToken = token
	return token, nil
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
