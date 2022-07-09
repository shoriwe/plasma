package ast3

type (
	Statement interface {
		S3()
	}
	Module struct {
		Statement
		Name *Identifier
		Body []Node
	}
	Assignment struct {
		Statement
		Left  Assignable
		Right Expression
	}
	Label struct {
		Statement
		Code int
	}
	Jump struct {
		Statement
		Target *Label
	}
	IfJump struct {
		Statement
		Condition Expression
		Target    *Label
	}
	Function struct {
		Statement
		Name      *Identifier
		Arguments []*Identifier
		Body      []Node
	}
	Return struct {
		Statement
		Result Expression
	}
	Require struct {
		Statement
		X Expression
	}

	Delete struct {
		Statement
		X Assignable
	}
	Defer struct {
		Statement
		X Expression
	}

	Block struct {
		Statement
		Body []Node
	}
)
