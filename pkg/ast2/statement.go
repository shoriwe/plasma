package ast2

const (
	HasNextString = "has_next"
	NextString    = "next"
)

type (
	AssignmentOperator int
	Statement          interface {
		Node
		S2()
	}
	Assignment struct {
		Statement
		Left  Assignable
		Right Expression
	}
	DoWhile struct {
		Statement
		Body      []Node
		Condition Expression
	}
	While struct {
		Statement
		Condition Expression
		Body      []Node
	}
	If struct {
		Statement
		SwitchSetup *Assignment
		Condition   Expression
		Body        []Node
		Else        []Node
	}
	Module struct {
		Statement
		Name *Identifier
		Body []Node
	}
	FunctionDefinition struct {
		Statement
		Name      *Identifier
		Arguments []*Identifier
		Body      []Node
	}
	GeneratorDefinition struct {
		Statement
		Name      *Identifier
		Arguments []*Identifier
		Body      []Node
	}
	Class struct {
		Statement
		Name  *Identifier
		Bases []Expression
		Body  []Node
	}
	ExceptBlock struct {
		Targets     []Expression
		CaptureName *Identifier
		Body        []Node
	}

	Try struct {
		Statement
		Body         []Node
		ExceptBlocks []*ExceptBlock
		Else         []Node
		Finally      []Node
	}
	Return struct {
		Statement
		Result Expression
	}
	Yield struct {
		Statement
		Result Expression
	}
	Continue struct {
		Statement
	}
	Break struct {
		Statement
	}
	Pass struct {
		Statement
	}
	Raise struct {
		Statement
		X Expression
	}
	Require struct {
		Statement
		X Expression
	}

	Delete struct {
		Statement
		X Assignable
	}

	Block struct {
		Statement
		Body []Node
	}
)
