package ast3

const (
	Not                = "not"
	Positive           = "positive"
	Negative           = "negative"
	NegateBits         = "negate_its"
	And                = "and"
	Or                 = "or"
	Xor                = "xor"
	In                 = "in"
	Is                 = "is"
	Implements         = "implements"
	Equals             = "equals"
	NotEqual           = "not_equal"
	GreaterThan        = "greater_than"
	GreaterOrEqualThan = "greater_or_equal_than"
	LessThan           = "less_than"
	LessOrEqualThan    = "less_or_equal_than"
	BitwiseOr          = "bitwise_or"
	BitwiseXor         = "bitwise_xor"
	BitwiseAnd         = "bitwise_and"
	BitwiseLeft        = "bitwise_left"
	BitwiseRight       = "bitwise_right"
	Add                = "add"
	Sub                = "sub"
	Mul                = "mul"
	Div                = "div"
	FloorDiv           = "floor_div"
	Modulus            = "mod"
	PowerOf            = "pow"
)

type (
	Expression interface {
		Node
		E3()
	}
	Assignable interface {
		Expression
		A2()
	}
	Call struct {
		Expression
		Function  Expression
		Arguments []Expression
	}
	IfOneLiner struct {
		Expression
		Condition, Result, Else Expression
	}

	Array struct {
		Expression
		Values []Expression
	}

	Tuple struct {
		Expression
		Values []Expression
	}

	KeyValue struct {
		Key, Value Expression
	}

	Hash struct {
		Expression
		Values []*KeyValue
	}

	Identifier struct {
		Assignable
		Symbol string
	}

	Integer struct {
		Expression
		Value int64
	}

	Float struct {
		Expression
		Value float64
	}

	String struct {
		Expression
		Contents []byte
	}

	Bytes struct {
		Expression
		Contents []byte
	}

	True struct {
		Expression
	}

	False struct {
		Expression
	}

	None struct {
		Expression
	}

	Selector struct {
		Assignable
		X          Expression
		Identifier *Identifier
	}

	Index struct {
		Assignable
		Source Expression
		Index  Expression
	}

	Super struct {
		Expression
		X Expression
	}
)
