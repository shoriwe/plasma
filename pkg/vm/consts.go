package vm

const (
	XXPrime1 uint64 = 11400714785074694791
	XXPrime2 uint64 = 14029467366897019727
	XXPrime5 uint64 = 2870177450012600261
)

const (
	TypeName          = "Type"
	ObjectName        = "Object"
	FunctionName      = "Function"
	StringName        = "String"
	BoolName          = "Bool"
	TrueName          = "True"
	FalseName         = "False"
	TupleName         = "Tuple"
	IntegerName       = "Integer"
	FloatName         = "Float"
	ArrayName         = "Array"
	NoneName          = "NoneType"
	BytesName         = "Bytes"
	HashName          = "Hash"
	IteratorName      = "Iterator"
	ModuleName        = "Module"
	None              = "None"
	CallableName      = "Callable"
	Source            = "0xFFFFFF"
	TemporalVariable1 = "0xAAAAAA"
	TemporalVariable2 = "0xBBBBBB"
	JunkVariable      = "0N-JUNK-VARIABLE"
)

const (
	Self                    = "self"
	Initialize              = "Initialize"
	NegBits                 = "NegBits"
	Negate                  = "Negate"
	Negative                = "Negative"
	Add                     = "Add"
	RightAdd                = "RightAdd"
	Sub                     = "Sub"
	RightSub                = "RightSub"
	Mul                     = "Mul"
	RightMul                = "RightMul"
	Div                     = "Div"
	RightDiv                = "RightDiv"
	FloorDiv                = "FloorDiv"
	RightFloorDiv           = "RightFloorDiv"
	Mod                     = "Mod"
	RightMod                = "RightMod"
	Pow                     = "Pow"
	RightPow                = "RightPow"
	BitXor                  = "BitXor"
	RightBitXor             = "RightBitXor"
	BitAnd                  = "BitAnd"
	RightBitAnd             = "RightBitAnd"
	BitOr                   = "BitOr"
	RightBitOr              = "RightBitOr"
	BitLeft                 = "BitLeft"
	RightBitLeft            = "RightBitLeft"
	BitRight                = "BitRight"
	RightBitRight           = "RightBitRight"
	And                     = "And"
	RightAnd                = "RightAnd"
	Or                      = "Or"
	RightOr                 = "RightOr"
	Xor                     = "Xor"
	RightXor                = "RightXor"
	Equals                  = "Equals"
	RightEquals             = "RightEquals"
	NotEquals               = "NotEquals"
	RightNotEquals          = "RightNotEquals"
	GreaterThan             = "GreaterThan"
	RightGreaterThan        = "RightGreaterThan"
	LessThan                = "LessThan"
	RightLessThan           = "RightLessThan"
	GreaterThanOrEqual      = "GreaterThanOrEqual"
	RightGreaterThanOrEqual = "RightGreaterThanOrEqual"
	LessThanOrEqual         = "LessThanOrEqual"
	RightLessThanOrEqual    = "RightLessThanOrEqual"
	Contains                = "Contains"
	RightContains           = "RightContains"
	Hash                    = "Hash"
	Copy                    = "Copy"
	Index                   = "Index"
	Assign                  = "Assign"
	Call                    = "Call"
	Iter                    = "Iter"
	HasNext                 = "HasNext"
	Next                    = "Next"
	Class                   = "Class"
	SubClasses              = "SubClasses"
	ToInteger               = "ToInteger"
	ToFloat                 = "ToFloat"
	ToString                = "ToString"
	ToBool                  = "ToBool"
	ToArray                 = "ToArray"
	ToTuple                 = "ToTuple"
	GetInteger              = "GetInteger"
	GetBool                 = "GetBool"
	GetBytes                = "GetBytes"
	GetString               = "GetString"
	GetFloat                = "GetFloat"
	GetContent              = "GetContent"
	GetKeyValues            = "GetKeyValues"
	GetLength               = "GetLength"
	SetBool                 = "SetBool"
	SetBytes                = "SetBytes"
	SetString               = "SetString"
	SetInteger              = "SetInteger"
	SetFloat                = "SetFloat"
	SetContent              = "SetContent"
	SetKeyValues            = "SetKeyValues"
	SetLength               = "SetLength"
)

const (
	Redo = iota
	Break
	NoAction
)

type ForLoopSettings struct {
	BodyLength int
	Receivers  []string
}

type IfInformation struct {
	Condition  []Code
	Body       []Code
	ElifBlocks []*IfInformation
	Else       []Code
}

type ExceptBlock struct {
	TargetErrors [][]Code
	Receiver     string
	Body         []Code
}

type TryInformation struct {
	Body         []Code
	ExceptBlocks []*ExceptBlock
	Else         []Code
	Finally      []Code
}
