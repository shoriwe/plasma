package lexer

func (lexer *Lexer) detectKindAndDirectValue() (Kind, DirectValue) {
	s := lexer.currentToken.String()
	switch s {
	case PassString:
		return Keyboard, Pass
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
	case ContinueString:
		return Keyboard, Continue
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
	case InterfaceString:
		return Keyboard, Interface
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
	case AsString:
		return Keyboard, As
	case RaiseString:
		return Keyboard, Raise
	case BEGINString:
		return Keyboard, BEGIN
	case ENDString:
		return Keyboard, END
	case NotString: // Unary operator
		return Operator, Not
	case TrueString:
		return Boolean, True
	case FalseString:
		return Boolean, False
	case NoneString:
		return NoneType, None
	case ContextString:
		return Keyboard, Context
	default:
		if identifierCheck.MatchString(s) {
			return IdentifierKind, InvalidDirectValue
		} else if junkKindCheck.MatchString(s) {
			return JunkKind, InvalidDirectValue
		}
	}
	return Unknown, InvalidDirectValue
}
