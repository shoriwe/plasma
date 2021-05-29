package vm

const (
	// Literal initialization
	NewStringOP uint8 = iota
	NewBytesOP
	NewIntegerOP
	NewFloatOP
	NewTrueBoolOP
	NewFalseBoolOP
	GetNoneOP

	// Composite creation
	NewTupleOP
	NewArrayOP
	NewHashOP

	// Unary Expressions
	NegateBitsOP
	BoolNegateOP

	// Binary Expressions
	AddOP
	SubOP
	MulOP
	DivOP
	ModOP
	PowOP
	BitXorOP
	BitAndOP
	BitOrOP
	BitLeftOP
	BitRightOP
	AndOP
	OrOP
	XorOP
	EqualsOP
	NotEqualsOP
	GreaterThanOP
	LessThanOP
	GreaterThanOrEqualOP
	LessThanOrEqualOP

	NoOP
	PushOP
	PushN_OP
	CopyOP
	CallOP
	GetOP
	GetFromOP
	ReturnOP
)
