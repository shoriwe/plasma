package lexer

import (
	"github.com/shoriwe/gplasma/pkg/errors"
	"github.com/shoriwe/gplasma/pkg/reader"
	"regexp"
)

var (
	identifierCheck = regexp.MustCompile("(?m)^[a-zA-Z_]+[a-zA-Z0-9_]*$")
	junkKindCheck   = regexp.MustCompile("(?m)^\\00+$")
)

type Lexer struct {
	currentToken *Token
	lastToken    *Token
	line         int
	reader       reader.Reader
	complete     bool
}

func (lexer *Lexer) HasNext() bool {
	return !lexer.complete
}

func (lexer *Lexer) next() (*Token, *errors.Error) {
	lexer.currentToken = &Token{
		Contents:    nil,
		DirectValue: InvalidDirectValue,
		Kind:        EOF,
		Line:        lexer.line,
		Index:       lexer.reader.Index(),
	}
	if !lexer.reader.HasNext() {
		lexer.complete = true
		return lexer.currentToken, nil
	}
	var tokenizingError *errors.Error
	char := lexer.reader.Char()
	lexer.reader.Next()
	switch char {
	case '\r':
		if lexer.reader.Char() != NewLineChar {
			return nil, errors.New(1, "invalid CRLF", "ssdsdsd")
		}
		lexer.reader.Next()
		lexer.line++
		lexer.currentToken.Contents = []rune{char, '\n'}
		lexer.currentToken.Kind = Separator
		lexer.currentToken.DirectValue = NewLine
	case NewLineChar:
		lexer.line++
		lexer.currentToken.Contents = []rune{char}
		lexer.currentToken.Kind = Separator
		lexer.currentToken.DirectValue = NewLine
	case SemiColonChar:
		lexer.currentToken.Contents = []rune{char}
		lexer.currentToken.Kind = Separator
		lexer.currentToken.DirectValue = SemiColon
	case ColonChar:
		lexer.currentToken.Contents = []rune{char}
		lexer.currentToken.DirectValue = Colon
		lexer.currentToken.Kind = Punctuation
	case CommaChar:
		lexer.currentToken.Contents = []rune{char}
		lexer.currentToken.DirectValue = Comma
		lexer.currentToken.Kind = Punctuation
	case OpenParenthesesChar:
		lexer.currentToken.Contents = []rune{char}
		lexer.currentToken.DirectValue = OpenParentheses
		lexer.currentToken.Kind = Punctuation
	case CloseParenthesesChar:
		lexer.currentToken.Contents = []rune{char}
		lexer.currentToken.DirectValue = CloseParentheses
		lexer.currentToken.Kind = Punctuation
	case OpenSquareBracketChar:
		lexer.currentToken.Contents = []rune{char}
		lexer.currentToken.DirectValue = OpenSquareBracket
		lexer.currentToken.Kind = Punctuation
	case CloseSquareBracketChar:
		lexer.currentToken.Contents = []rune{char}
		lexer.currentToken.DirectValue = CloseSquareBracket
		lexer.currentToken.Kind = Punctuation
	case OpenBraceChar:
		lexer.currentToken.Contents = []rune{char}
		lexer.currentToken.DirectValue = OpenBrace
		lexer.currentToken.Kind = Punctuation
	case CloseBraceChar:
		lexer.currentToken.Contents = []rune{char}
		lexer.currentToken.DirectValue = CloseBrace
		lexer.currentToken.Kind = Punctuation
	case DollarSignChar:
		lexer.currentToken.Contents = []rune{char}
		lexer.currentToken.DirectValue = DollarSign
		lexer.currentToken.Kind = Punctuation
	case DotChar:
		lexer.currentToken.Contents = []rune{char}
		lexer.currentToken.DirectValue = Dot
		lexer.currentToken.Kind = Punctuation
	case WhiteSpaceChar:
		lexer.currentToken.DirectValue = Blank
		lexer.currentToken.Kind = Whitespace
		lexer.currentToken.Contents = []rune{char}
	case TabChar:
		lexer.currentToken.DirectValue = Blank
		lexer.currentToken.Kind = Whitespace
		lexer.currentToken.Contents = []rune{char}
	case CommentChar:
		lexer.currentToken.Contents, lexer.currentToken.Kind, tokenizingError = lexer.tokenizeComment()
		lexer.currentToken.Contents = append([]rune{'#'}, lexer.currentToken.Contents...)
	case '\'', '"': // String1
		lexer.currentToken.Contents, lexer.currentToken.Kind, lexer.currentToken.DirectValue, tokenizingError = lexer.tokenizeStringLikeExpressions(char)
	case '`':
		lexer.currentToken.Contents, lexer.currentToken.Kind, lexer.currentToken.DirectValue, tokenizingError = lexer.tokenizeStringLikeExpressions(char)
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		lexer.currentToken.Contents, lexer.currentToken.Kind, lexer.currentToken.DirectValue, tokenizingError = lexer.tokenizeNumeric(char)
	case StarChar:
		lexer.currentToken.Contents, lexer.currentToken.Kind, lexer.currentToken.DirectValue = lexer.tokenizeRepeatableOperator(char, Star, Operator, PowerOf, Operator, StarAssign, Assignment, PowerOfAssign, Assignment)
	case DivChar:
		lexer.currentToken.Contents, lexer.currentToken.Kind, lexer.currentToken.DirectValue = lexer.tokenizeRepeatableOperator(char, Div, Operator, FloorDiv, Operator, DivAssign, Assignment, FloorDivAssign, Assignment)
	case LessThanChar:
		lexer.currentToken.Contents, lexer.currentToken.Kind, lexer.currentToken.DirectValue = lexer.tokenizeRepeatableOperator(char, LessThan, Comparator, BitwiseLeft, Operator, LessOrEqualThan, Comparator, BitwiseLeftAssign, Assignment)
	case GreatThanChar:
		lexer.currentToken.Contents, lexer.currentToken.Kind, lexer.currentToken.DirectValue = lexer.tokenizeRepeatableOperator(char, GreaterThan, Comparator, BitwiseRight, Operator, GreaterOrEqualThan, Comparator, BitwiseRightAssign, Assignment)
	case AddChar:
		lexer.currentToken.Contents, lexer.currentToken.Kind, lexer.currentToken.DirectValue = lexer.tokenizeSingleOperator(char, Add, Operator, AddAssign, Assignment)
	case SubChar:
		lexer.currentToken.Contents, lexer.currentToken.Kind, lexer.currentToken.DirectValue = lexer.tokenizeSingleOperator(char, Sub, Operator, SubAssign, Assignment)
	case ModulusChar:
		lexer.currentToken.Contents, lexer.currentToken.Kind, lexer.currentToken.DirectValue = lexer.tokenizeSingleOperator(char, Modulus, Operator, ModulusAssign, Assignment)
	case BitwiseXorChar:
		lexer.currentToken.Contents, lexer.currentToken.Kind, lexer.currentToken.DirectValue = lexer.tokenizeSingleOperator(char, BitwiseXor, Operator, BitwiseXorAssign, Assignment)
	case BitWiseAndChar:
		lexer.currentToken.Contents, lexer.currentToken.Kind, lexer.currentToken.DirectValue = lexer.tokenizeSingleOperator(char, BitWiseAnd, Operator, BitWiseAndAssign, Assignment)
	case BitwiseOrChar:
		lexer.currentToken.Contents, lexer.currentToken.Kind, lexer.currentToken.DirectValue = lexer.tokenizeSingleOperator(char, BitwiseOr, Operator, BitwiseOrAssign, Assignment)
	case SignNotChar:
		lexer.currentToken.Contents, lexer.currentToken.Kind, lexer.currentToken.DirectValue = lexer.tokenizeSingleOperator(char, SignNot, Operator, NotEqual, Comparator)
	case NegateBitsChar:
		lexer.currentToken.Contents, lexer.currentToken.Kind, lexer.currentToken.DirectValue = lexer.tokenizeSingleOperator(char, NegateBits, Operator, NegateBitsAssign, Assignment)
	case EqualsChar:
		lexer.currentToken.Contents, lexer.currentToken.Kind, lexer.currentToken.DirectValue = lexer.tokenizeSingleOperator(char, Assign, Assignment, Equals, Comparator)
	case BackSlashChar:
		lexer.currentToken.Contents = []rune{char}
		if lexer.reader.HasNext() {
			nextChar := lexer.reader.Char()
			if nextChar != '\n' {
				return nil, errors.New(lexer.line, "line escape not followed by a new line", errors.LexingError)
			}
			lexer.currentToken.Contents = append(lexer.currentToken.Contents, '\n')
			lexer.reader.Next()
		}
		lexer.currentToken.Kind = PendingEscape
	default:
		if char == 'b' {
			if lexer.reader.HasNext() {
				nextChar := lexer.reader.Char()
				if nextChar == '\'' || nextChar == '"' {
					var byteStringPart []rune
					lexer.reader.Next()
					byteStringPart, lexer.currentToken.Kind, lexer.currentToken.DirectValue, tokenizingError = lexer.tokenizeStringLikeExpressions(nextChar)
					lexer.currentToken.Contents = append([]rune{char}, byteStringPart...)
					if lexer.currentToken.DirectValue != InvalidDirectValue {
						lexer.currentToken.DirectValue = ByteString
					}
					break
				}
			}
		}
		lexer.currentToken.Contents, lexer.currentToken.Kind, lexer.currentToken.DirectValue,
			tokenizingError = lexer.tokenizeWord(char)
	}
	return lexer.currentToken, tokenizingError
}

/*
	This function will yield just the necessary token, this means not repeated separators
*/
func (lexer *Lexer) Next() (*Token, *errors.Error) {
	token, tokenizingError := lexer.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if token.Kind == JunkKind {
		return lexer.Next()
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

func NewLexer(codeReader reader.Reader) *Lexer {
	return &Lexer{
		lastToken: nil,
		line:      1,
		reader:    codeReader,
		complete:  false,
	}
}
