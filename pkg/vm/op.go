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
	NegativeOP

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
	// Other expressions
	GetIdentifierOP
	IndexOP
	SelectNameFromObjectOP
	MethodInvocationOP

	// Assign Statement
	AssignIdentifierOP
	AssignSelectorOP
	AssignIndexOP
	IfJumpOP
	UnlessJumpOP
	BreakOP
	RedoOP
	ContinueOP

	ReturnOP

	// Special Instructions
	LoadFunctionArgumentsOP
	NewFunctionOP
	JumpOP
	PopOP
	NOP
	SetupForLoopOP
	HasNextOP
	UnpackReceiversPopOP
	UnpackReceiversPeekOP
	PopIterOP
	NewIteratorOP

	// Try-Except related
	SetupTryBlockOP
	SetupTryExceptBlockOP
	SetupTryElseBlockOP
	SetupTryFinallyBlockOP
	ExitTryBlockOP

	NewModuleOP
	NewClassOP
	NewClassFunctionOP

	RaiseOP
)
