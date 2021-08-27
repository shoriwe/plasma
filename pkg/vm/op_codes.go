package vm

const (
	ReturnState uint8 = iota
	BreakState
	RedoState
	ContinueState
	NoState
)

const (
	NewStringOP uint8 = iota
	NewBytesOP
	NewIntegerOP
	NewFloatOP
	GetTrueOP
	GetFalseOP
	NewLambdaFunctionOP
	GetNoneOP
	NewTupleOP
	NewArrayOP
	NewHashOP

	UnaryOP
	NegateBitsOP
	BoolNegateOP
	NegativeOP

	BinaryOP
	AddOP
	SubOP
	MulOP
	DivOP
	FloorDivOP
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
	ContainsOP

	GetIdentifierOP
	IndexOP
	SelectNameFromObjectOP
	MethodInvocationOP
	AssignIdentifierOP
	AssignSelectorOP
	AssignIndexOP
	IfJumpOP
	UnlessJumpOP
	BreakOP
	RedoOP
	ContinueOP
	ReturnOP
	ForLoopOP
	LoadFunctionArgumentsOP
	NewFunctionOP
	JumpOP
	PushOP
	PopOP
	NOP
	NewGeneratorOP
	TryOP
	NewModuleOP
	NewClassOP
	NewClassFunctionOP
	RaiseOP
)

var unaryInstructionsFunctions = map[uint8]string{
	NegateBitsOP: NegateBits,
	BoolNegateOP: Negate,
	NegativeOP:   Negative,
}

var binaryInstructionsFunctions = map[uint8][2]string{
	AddOP:                {Add, RightAdd},
	SubOP:                {Sub, RightSub},
	MulOP:                {Mul, RightMul},
	DivOP:                {Div, RightDiv},
	FloorDivOP:           {FloorDiv, RightFloorDiv},
	ModOP:                {Mod, RightMod},
	PowOP:                {Pow, RightPow},
	BitXorOP:             {BitXor, RightBitXor},
	BitAndOP:             {BitAnd, RightBitAnd},
	BitOrOP:              {BitOr, RightBitOr},
	BitLeftOP:            {BitLeft, RightBitLeft},
	BitRightOP:           {BitRight, RightBitRight},
	AndOP:                {And, RightAnd},
	OrOP:                 {Or, RightOr},
	XorOP:                {Xor, RightXor},
	EqualsOP:             {Equals, RightEquals},
	NotEqualsOP:          {NotEquals, RightNotEquals},
	GreaterThanOP:        {GreaterThan, RightGreaterThan},
	LessThanOP:           {LessThan, RightLessThan},
	GreaterThanOrEqualOP: {GreaterThanOrEqual, RightGreaterThanOrEqual},
	LessThanOrEqualOP:    {LessThanOrEqual, RightLessThanOrEqual},
	ContainsOP:           {"029p3847980479087437891734", Contains},
}

var instructionNames = map[uint8]string{
	NewStringOP:         "NewStringOP",
	NewBytesOP:          "NewBytesOP",
	NewIntegerOP:        "NewIntegerOP",
	NewFloatOP:          "NewFloatOP",
	GetTrueOP:           "GetTrueOP",
	GetFalseOP:          "GetFalseOP",
	NewLambdaFunctionOP: "NewLambdaFunctionOP",
	GetNoneOP:           "GetNoneOP",
	NewTupleOP:          "NewTupleOP",
	NewArrayOP:          "NewArrayOP",
	NewHashOP:           "NewHashOP",

	UnaryOP:      "UnaryOP",
	NegateBitsOP: "NegateBitsOP",
	BoolNegateOP: "BoolNegateOP",
	NegativeOP:   "NegativeOP",

	BinaryOP:             "BinaryOP",
	AddOP:                "AddOP",
	SubOP:                "SubOP",
	MulOP:                "MulOP",
	DivOP:                "DivOP",
	FloorDivOP:           "FloorDivOP",
	ModOP:                "ModOP",
	PowOP:                "PowOP",
	BitXorOP:             "BitXorOP",
	BitAndOP:             "BitAndOP",
	BitOrOP:              "BitOrOP",
	BitLeftOP:            "BitLeftOP",
	BitRightOP:           "BitRightOP",
	AndOP:                "AndOP",
	OrOP:                 "OrOP",
	XorOP:                "XorOP",
	EqualsOP:             "EqualsOP",
	NotEqualsOP:          "NotEqualsOP",
	GreaterThanOP:        "GreaterThanOP",
	LessThanOP:           "LessThanOP",
	GreaterThanOrEqualOP: "GreaterThanOrEqualOP",
	LessThanOrEqualOP:    "LessThanOrEqualOP",
	ContainsOP:           "ContainsOP",

	GetIdentifierOP:         "GetIdentifierOP",
	IndexOP:                 "IndexOP",
	SelectNameFromObjectOP:  "SelectNameFromObjectOP",
	MethodInvocationOP:      "MethodInvocationOP",
	AssignIdentifierOP:      "AssignIdentifierOP",
	AssignSelectorOP:        "AssignSelectorOP",
	AssignIndexOP:           "AssignIndexOP",
	IfJumpOP:                "IfJumpOP",
	UnlessJumpOP:            "UnlessJumpOP",
	BreakOP:                 "BreakOP",
	RedoOP:                  "RedoOP",
	ContinueOP:              "ContinueOP",
	ReturnOP:                "ReturnOP",
	ForLoopOP:               "ForLoopOP",
	LoadFunctionArgumentsOP: "LoadFunctionArgumentsOP",
	NewFunctionOP:           "NewFunctionOP",
	JumpOP:                  "JumpOP",
	PushOP:                  "PushOP",
	PopOP:                   "PopOP",
	NOP:                     "NOP",
	NewGeneratorOP:          "NewGeneratorOP",
	TryOP:                   "TryOP",
	NewModuleOP:             "NewModuleOP",
	NewClassOP:              "NewClassOP",
	NewClassFunctionOP:      "NewClassFunctionOP",
	RaiseOP:                 "RaiseOP",
}
