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
		lexer.currentToken.append(char, '\n')
		lexer.currentToken.Kind = Separator
		lexer.currentToken.DirectValue = NewLine
	case NewLineChar:
		lexer.line++
		lexer.currentToken.append(char)
		lexer.currentToken.Kind = Separator
		lexer.currentToken.DirectValue = NewLine
	case SemiColonChar:
		lexer.currentToken.append(char)
		lexer.currentToken.Kind = Separator
		lexer.currentToken.DirectValue = SemiColon
	case ColonChar:
		lexer.currentToken.append(char)
		lexer.currentToken.DirectValue = Colon
		lexer.currentToken.Kind = Punctuation
	case CommaChar:
		lexer.currentToken.append(char)
		lexer.currentToken.DirectValue = Comma
		lexer.currentToken.Kind = Punctuation
	case OpenParenthesesChar:
		lexer.currentToken.append(char)
		lexer.currentToken.DirectValue = OpenParentheses
		lexer.currentToken.Kind = Punctuation
	case CloseParenthesesChar:
		lexer.currentToken.append(char)
		lexer.currentToken.DirectValue = CloseParentheses
		lexer.currentToken.Kind = Punctuation
	case OpenSquareBracketChar:
		lexer.currentToken.append(char)
		lexer.currentToken.DirectValue = OpenSquareBracket
		lexer.currentToken.Kind = Punctuation
	case CloseSquareBracketChar:
		lexer.currentToken.append(char)
		lexer.currentToken.DirectValue = CloseSquareBracket
		lexer.currentToken.Kind = Punctuation
	case OpenBraceChar:
		lexer.currentToken.append(char)
		lexer.currentToken.DirectValue = OpenBrace
		lexer.currentToken.Kind = Punctuation
	case CloseBraceChar:
		lexer.currentToken.append(char)
		lexer.currentToken.DirectValue = CloseBrace
		lexer.currentToken.Kind = Punctuation
	case DollarSignChar:
		lexer.currentToken.append(char)
		lexer.currentToken.DirectValue = DollarSign
		lexer.currentToken.Kind = Punctuation
	case DotChar:
		lexer.currentToken.append(char)
		lexer.currentToken.DirectValue = Dot
		lexer.currentToken.Kind = Punctuation
	case WhiteSpaceChar:
		lexer.currentToken.append(char)
		lexer.currentToken.DirectValue = Blank
		lexer.currentToken.Kind = Whitespace
	case TabChar:
		lexer.currentToken.append(char)
		lexer.currentToken.DirectValue = Blank
		lexer.currentToken.Kind = Whitespace
	case CommentChar:
		lexer.currentToken.append(char)
		lexer.tokenizeComment()
	case '\'', '"', '`':
		lexer.currentToken.append(char)
		tokenizingError = lexer.tokenizeStringLikeExpressions(char)
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		lexer.currentToken.append(char)
		tokenizingError = lexer.tokenizeNumeric()
	case StarChar:
		lexer.currentToken.append(char)
		lexer.tokenizeRepeatableOperator(Star, Operator, PowerOf, Operator, StarAssign, Assignment, PowerOfAssign, Assignment)
	case DivChar:
		lexer.currentToken.append(char)
		lexer.tokenizeRepeatableOperator(Div, Operator, FloorDiv, Operator, DivAssign, Assignment, FloorDivAssign, Assignment)
	case LessThanChar:
		lexer.currentToken.append(char)
		lexer.tokenizeRepeatableOperator(LessThan, Comparator, BitwiseLeft, Operator, LessOrEqualThan, Comparator, BitwiseLeftAssign, Assignment)
	case GreatThanChar:
		lexer.currentToken.append(char)
		lexer.tokenizeRepeatableOperator(GreaterThan, Comparator, BitwiseRight, Operator, GreaterOrEqualThan, Comparator, BitwiseRightAssign, Assignment)
	case AddChar:
		lexer.currentToken.append(char)
		lexer.tokenizeSingleOperator(Add, Operator, AddAssign, Assignment)
	case SubChar:
		lexer.currentToken.append(char)
		lexer.tokenizeSingleOperator(Sub, Operator, SubAssign, Assignment)
	case ModulusChar:
		lexer.currentToken.append(char)
		lexer.tokenizeSingleOperator(Modulus, Operator, ModulusAssign, Assignment)
	case BitwiseXorChar:
		lexer.currentToken.append(char)
		lexer.tokenizeSingleOperator(BitwiseXor, Operator, BitwiseXorAssign, Assignment)
	case BitWiseAndChar:
		lexer.currentToken.append(char)
		lexer.tokenizeSingleOperator(BitWiseAnd, Operator, BitWiseAndAssign, Assignment)
	case BitwiseOrChar:
		lexer.currentToken.append(char)
		lexer.tokenizeSingleOperator(BitwiseOr, Operator, BitwiseOrAssign, Assignment)
	case SignNotChar:
		lexer.currentToken.append(char)
		lexer.tokenizeSingleOperator(SignNot, Operator, NotEqual, Comparator)
	case NegateBitsChar:
		lexer.currentToken.append(char)
		lexer.tokenizeSingleOperator(NegateBits, Operator, NegateBitsAssign, Assignment)
	case EqualsChar:
		lexer.currentToken.append(char)
		lexer.tokenizeSingleOperator(Assign, Assignment, Equals, Comparator)
	case BackSlashChar:
		lexer.currentToken.append(char)
		if !lexer.reader.HasNext() {
			return nil, errors.New(lexer.line, "line escape not followed by a new line", errors.LexingError)
		}
		nextChar := lexer.reader.Char()
		if nextChar != '\n' {
			return nil, errors.New(lexer.line, "line escape not followed by a new line", errors.LexingError)
		}
		lexer.currentToken.append('\n')
		lexer.reader.Next()

	default:
		if char != 'b' || !lexer.reader.HasNext() {
			lexer.currentToken.Contents, lexer.currentToken.Kind, lexer.currentToken.DirectValue,
				tokenizingError = lexer.tokenizeWord(char)
			break
		}
		nextChar := lexer.reader.Char()
		if nextChar != '\'' && nextChar != '"' {
			lexer.currentToken.Contents, lexer.currentToken.Kind, lexer.currentToken.DirectValue,
				tokenizingError = lexer.tokenizeWord(char)
			break
		}
		lexer.reader.Next()
		lexer.currentToken.append('b', nextChar)
		tokenizingError = lexer.tokenizeStringLikeExpressions(nextChar)
		if lexer.currentToken.DirectValue != InvalidDirectValue {
			lexer.currentToken.DirectValue = ByteString
		}
		break
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
